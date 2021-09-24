package tstbuild

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/trustelem/esbuild/pkg/api"
	"github.com/trustelem/esbuild/pkg/cli"
)

func checkEnv(t *testing.T) {
	t.Helper()
	tag := os.Getenv("ESBUILD_TESTCLI")
	if tag == "" || tag == "0" {
		t.Skip("Set ESBUILD_TESTCLI=1 to run this test")
	}
}

func toJson(v interface{}, indent bool) string {
	var bs []byte
	var err error
	if indent {
		bs, err = json.MarshalIndent(v, "", "  ")
	} else {
		bs, err = json.Marshal(v)
	}

	if err != nil {
		return fmt.Sprintf(`<error:%s>`, err)
	}
	return string(bs)
}

func TestApiBuildNamedFile(t *testing.T) {
	checkEnv(t)

	const output = "api.js"
	if _, err := os.Stat(output); err == nil {
		os.Remove(output)
	}

	opts, _, _, _ := cli.UnitTestParseOptions([]string{"--bundle", "--outfile=" + output, "app.jsx"})

	result := api.Build(*opts)
	if len(result.Errors) > 0 {
		t.Fatalf("build failed:\n%s", toJson(result.Errors, true))
	}
}

func TestApiBuildMinifiedNamedFile(t *testing.T) {
	checkEnv(t)

	const output = "api.min.js"
	if _, err := os.Stat(output); err == nil {
		os.Remove(output)
	}

	opts, _, _, _ := cli.UnitTestParseOptions([]string{"--bundle", "--minify", "--outfile=" + output, "app.jsx"})

	result := api.Build(*opts)
	if len(result.Errors) > 0 {
		t.Fatalf("build failed:\n%s", toJson(result.Errors, true))
	}
}
