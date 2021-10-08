package tstbuild

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/trustelem/esbuild/pkg/api"
	"github.com/trustelem/esbuild/pkg/cli"
	"github.com/trustelem/esbuild/pkg/zipfs"
)

func TestApiZipFS(t *testing.T) {
	skipTestIfEnvNotSet(t)

	cwd, _ := os.Getwd()
	root := filepath.Join(cwd, "node_modules")

	zfs, err := zipfs.New("foo.zip", root, "node_modules/", cwd)
	if err != nil {
		t.Fatal(err)
	}

	const output = "zip.js"
	if _, err := os.Stat(output); err == nil {
		os.Remove(output)
	}

	opts, _, _, _ := cli.UnitTestParseOptions([]string{"--bundle", "--minify", "--outfile=" + output, "app.jsx"})
	opts.FS = zfs

	result := api.Build(*opts)
	if len(result.Errors) > 0 {
		t.Fatalf("build failed\n%s", toJson(result.Errors, true))
	}
}
