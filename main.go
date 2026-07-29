package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/engine"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/githubactions"

	_ "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules/azurerm"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "upgrader",
		Short: "Upgrade Terraform azurerm configuration across major versions",
	}
	root.AddCommand(updateCmd())
	root.AddCommand(validateCmd())
	return root
}

func updateCmd() *cobra.Command {
	var (
		dir        string
		dryRun     bool
		reportKind string
		recursive  bool
		format     bool
	)
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Apply azurerm v5 breaking-change transforms to a module directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			// --fmt=false preserves the original formatting on lines the upgrade
			// did not change; the engine option is phrased as the opt-out.
			code, err := runUpdate(cmd.OutOrStdout(), dir, dryRun, recursive, !format, reportKind)
			if err != nil {
				return err
			}
			os.Exit(code)
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "module directory to upgrade")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print a diff instead of writing files")
	cmd.Flags().StringVar(&reportKind, "report", "text", "report format: text|json")
	cmd.Flags().BoolVar(&recursive, "recursive", true,
		"recurse into nested module directories (matches validate); use --recursive=false for the top-level directory only")
	cmd.Flags().BoolVar(&format, "fmt", true,
		"canonically format changed files, like terraform fmt; use --fmt=false to keep original formatting on lines the upgrade did not change")
	return cmd
}

// runUpdate applies the v5 transforms to dir (recursively unless recursive is
// false, mirroring validate's discovery), renders the report to w, and returns
// the process exit code. It is split out of the cobra command so it is testable
// without the command's os.Exit.
func runUpdate(w io.Writer, dir string, dryRun, recursive, preserveFormatting bool, reportKind string) (exitCode int, err error) {
	if reportKind != "text" && reportKind != "json" {
		return 0, fmt.Errorf("invalid --report %q (want text|json)", reportKind)
	}
	res, err := engine.Run(engine.Options{Dir: dir, DryRun: dryRun, Recursive: recursive, PreserveFormatting: preserveFormatting})
	if err != nil {
		return 0, err
	}
	if dryRun {
		for path, d := range res.Diffs {
			fmt.Fprintf(w, "# %s\n%s\n", path, d)
		}
	}
	switch reportKind {
	case "json":
		if err := res.Report.RenderJSON(w); err != nil {
			return 0, err
		}
	default:
		if err := res.Report.RenderText(w); err != nil {
			return 0, err
		}
	}
	return res.Report.ExitCode(), nil
}

func validateCmd() *cobra.Command {
	var (
		reportKind    string
		githubActions bool
	)
	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Read-only check for whether Terraform configurations need azurerm v5 upgrades",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			needsUpgrade, err := runValidate(cmd.OutOrStdout(), path, reportKind, githubActions,
				githubactions.LookupEnv(os.Getenv), http.DefaultClient)
			if err != nil {
				return err
			}
			if needsUpgrade {
				os.Exit(1)
			}
			os.Exit(0)
			return nil
		},
	}
	cmd.Flags().StringVar(&reportKind, "report", "text", "report format: text|json")
	cmd.Flags().BoolVar(&githubActions, "github-actions", false,
		"additionally emit GitHub Actions channels (annotations, step summary, outputs, PR comment)")
	return cmd
}

// detect performs a read-only, recursive, dry-run engine pass over dir and
// reports whether any discovered module needs an azurerm v5 upgrade. It never
// writes to disk.
func detect(dir string) (needsUpgrade bool, res *engine.Result, err error) {
	res, err = engine.Run(engine.Options{Dir: dir, DryRun: true, Recursive: true})
	if err != nil {
		return false, nil, err
	}
	needsUpgrade = len(res.Diffs) > 0 || len(res.Report.Findings) > 0
	return needsUpgrade, res, nil
}

// runValidate performs a read-only, recursive, dry-run detection over dir and
// renders the report to w. It never writes to disk. When githubActions is set it
// additionally emits the GitHub Actions channels (annotations, step summary,
// outputs, PR comment). It returns whether any discovered module needs an
// azurerm v5 upgrade, which callers use to pick an exit code.
func runValidate(w io.Writer, dir, reportKind string, githubActions bool, env githubactions.Environment, client *http.Client) (needsUpgrade bool, err error) {
	if reportKind != "text" && reportKind != "json" {
		return false, fmt.Errorf("invalid --report %q (want text|json)", reportKind)
	}
	needsUpgrade, res, err := detect(dir)
	if err != nil {
		return false, err
	}

	switch reportKind {
	case "json":
		if err := res.Report.RenderJSONValidate(w, needsUpgrade); err != nil {
			return needsUpgrade, err
		}
	default:
		if err := res.Report.RenderText(w); err != nil {
			return needsUpgrade, err
		}
		if _, err := fmt.Fprintf(w, "needs upgrade: %t\n", needsUpgrade); err != nil {
			return needsUpgrade, err
		}
	}

	if githubActions {
		if err := emitGitHubActions(w, res, needsUpgrade, env, client); err != nil {
			return needsUpgrade, err
		}
	}
	return needsUpgrade, nil
}

// emitGitHubActions renders the shared markdown body and drives the GitHub
// Actions channels: inline annotations (to stdout), the step summary and action
// outputs (when their env paths are set), and a pull-request comment. An error
// signals an operational failure (e.g. a missing token in a PR context).
func emitGitHubActions(stdout io.Writer, res *engine.Result, needsUpgrade bool, env githubactions.Environment, client *http.Client) error {
	var md bytes.Buffer
	if err := res.Report.RenderMarkdown(&md, needsUpgrade); err != nil {
		return err
	}

	deps := githubactions.Deps{
		Env:        env,
		Stdout:     stdout,
		ReadFile:   os.ReadFile,
		HTTPClient: client,
	}
	if env.SummaryPath != "" {
		f, err := os.OpenFile(env.SummaryPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		deps.Summary = f
	}
	if env.OutputPath != "" {
		f, err := os.OpenFile(env.OutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		deps.Output = f
	}

	return githubactions.Emit(deps, res.Report.Findings, md.String(), needsUpgrade)
}
