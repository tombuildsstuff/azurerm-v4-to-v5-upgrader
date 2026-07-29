// Package githubactions renders azurerm v5 upgrade detection results through
// GitHub Actions channels: inline annotations, the step summary, action
// outputs, and a pull-request comment.
package githubactions

// Environment captures the GitHub Actions context, read from the standard
// GITHUB_* environment variables.
type Environment struct {
	IsActions   bool
	Repository  string // "owner/repo"
	EventName   string
	EventPath   string
	Token       string
	APIURL      string // defaults to https://api.github.com
	SummaryPath string // $GITHUB_STEP_SUMMARY
	OutputPath  string // $GITHUB_OUTPUT
}

// LookupEnv builds an Environment from getenv (normally os.Getenv). APIURL
// defaults to the public GitHub API when GITHUB_API_URL is unset.
func LookupEnv(getenv func(string) string) Environment {
	apiURL := getenv("GITHUB_API_URL")
	if apiURL == "" {
		apiURL = "https://api.github.com"
	}
	return Environment{
		IsActions:   getenv("GITHUB_ACTIONS") == "true",
		Repository:  getenv("GITHUB_REPOSITORY"),
		EventName:   getenv("GITHUB_EVENT_NAME"),
		EventPath:   getenv("GITHUB_EVENT_PATH"),
		Token:       getenv("GITHUB_TOKEN"),
		APIURL:      apiURL,
		SummaryPath: getenv("GITHUB_STEP_SUMMARY"),
		OutputPath:  getenv("GITHUB_OUTPUT"),
	}
}

// IsPullRequest reports whether the workflow was triggered by a pull request.
func (e Environment) IsPullRequest() bool {
	return e.EventName == "pull_request" || e.EventName == "pull_request_target"
}
