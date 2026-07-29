package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// File holds two parallel views of the same *.tf source: a write tree for
// surgical edits and a syntax body for analysis, plus the original bytes.
type File struct {
	Path     string
	Write    *hclwrite.File
	Syntax   *hclsyntax.Body
	Original []byte
}

// Load parses every *.tf file in dir (non-recursive) into paired views.
// It returns an error if any file fails to parse.
func Load(dir string) ([]*File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tf") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	var files []*File
	for _, name := range names {
		path := filepath.Join(dir, name)
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		f, err := parse(path, src)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func parse(path string, src []byte) (*File, error) {
	wf, wdiags := hclwrite.ParseConfig(src, path, hcl.InitialPos)
	if wdiags.HasErrors() {
		return nil, fmt.Errorf("parsing %s: %s", path, wdiags.Error())
	}
	sf, sdiags := hclsyntax.ParseConfig(src, path, hcl.InitialPos)
	if sdiags.HasErrors() {
		return nil, fmt.Errorf("parsing %s: %s", path, sdiags.Error())
	}
	body, ok := sf.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("parsing %s: unexpected body type", path)
	}
	return &File{Path: path, Write: wf, Syntax: body, Original: src}, nil
}

// LoadTree parses every *.tf file under root and all nested directories into
// paired views. It skips any directory (other than root itself) whose base name
// begins with "." - e.g. .terraform and .git - and returns files sorted by path.
func LoadTree(root string) ([]*File, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".tf") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", root, err)
	}
	sort.Strings(paths)

	var files []*File
	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		f, err := parse(path, src)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}
