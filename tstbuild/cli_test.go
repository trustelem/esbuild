package tstbuild

import (
	"os"
	"testing"

	"github.com/trustelem/esbuild/pkg/cli"
)

func TestBuildNamedFile(t *testing.T) {
	checkEnv(t)

	const output = "cli.js"
	if _, err := os.Stat(output); err == nil {
		os.Remove(output)
	}

	exitCode := cli.Run([]string{"--bundle", "--outfile=" + output, "app.jsx"})
	if exitCode != 0 {
		t.Fatal("build failed")
	}
}

func TestBuildMinifiedNamedFile(t *testing.T) {
	checkEnv(t)

	const output = "cli.min.js"
	if _, err := os.Stat(output); err == nil {
		os.Remove(output)
	}

	exitCode := cli.Run([]string{"--bundle", "--minify", "--outfile=" + output, "app.jsx"})
	if exitCode != 0 {
		t.Fatal("build failed")
	}
}
