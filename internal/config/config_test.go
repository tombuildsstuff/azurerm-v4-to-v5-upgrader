package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, dir, name, body string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644))
}

func TestLoadParsesTFFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.tf", "resource \"azurerm_resource_group\" \"a\" {\n  name = \"x\"\n}\n")
	writeFile(t, dir, "vars.tf", "variable \"y\" {}\n")
	writeFile(t, dir, "ignore.txt", "not terraform")

	files, err := Load(dir)
	require.NoError(t, err)
	require.Len(t, files, 2)
	require.Equal(t, "main.tf", filepath.Base(files[0].Path))
	require.Equal(t, "vars.tf", filepath.Base(files[1].Path))
	require.NotNil(t, files[0].Write)
	require.NotNil(t, files[0].Syntax)
	require.Equal(t, "resource \"azurerm_resource_group\" \"a\" {\n  name = \"x\"\n}\n", string(files[0].Original))
}

func TestLoadReturnsErrorOnInvalidHCL(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "bad.tf", "resource \"azurerm_x\" \"a\" {")
	_, err := Load(dir)
	require.Error(t, err)
}

func TestLoadTreeRecursesAndSkipsDotDirs(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(root, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	}
	write("a.tf", "x = 1\n")
	write(filepath.Join("sub", "b.tf"), "y = 2\n")
	write(filepath.Join(".terraform", "c.tf"), "z = 3\n")
	write(filepath.Join("sub", "notes.txt"), "ignored\n")

	files, err := LoadTree(root)
	require.NoError(t, err)

	var paths []string
	for _, f := range files {
		rel, relErr := filepath.Rel(root, f.Path)
		require.NoError(t, relErr)
		paths = append(paths, rel)
	}
	require.Equal(t, []string{"a.tf", filepath.Join("sub", "b.tf")}, paths)
}
