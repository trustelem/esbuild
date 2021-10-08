package api

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func skipTestIfEnvNotSet(t *testing.T) {
	t.Helper()
	tag := os.Getenv("ESBUILD_GEN_TEST")
	if tag == "" || tag == "0" {
		t.Skip("Set ESBUILD_GEN_TEST=1 to run this test")
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
	//skipTestIfEnvNotSet(t)

	{
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		err = os.Chdir("../../tstbuild")
		if err != nil {
			t.Fatal(err)
		}

		defer os.Chdir(cwd)
	}

	const output = "api.js"
	if _, err := os.Stat(output); err == nil {
		os.Remove(output)
	}

	opts := &BuildOptions{
		LogLimit:    10,
		LogLevel:    3,
		Bundle:      true,
		Outfile:     output,
		EntryPoints: []string{"app.jsx"},
		Write:       true,
	}

	result := Build(*opts)
	if len(result.Errors) > 0 {
		t.Fatalf("build failed:\n%s", toJson(result.Errors, true))
	}
}

func TestApiBuildMinifiedNamedFile(t *testing.T) {
	//skipTestIfEnvNotSet(t)

	const output = "api.min.js"
	if _, err := os.Stat(output); err == nil {
		os.Remove(output)
	}

	opts := &BuildOptions{
		LogLimit:          10,
		LogLevel:          3,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Bundle:            true,
		Outfile:           output,
		EntryPoints:       []string{"app.jsx"},
		Write:             true,
	}

	t.Logf("options: %s", toJson(opts, true))

	result := Build(*opts)
	if len(result.Errors) > 0 {
		t.Fatalf("build failed:\n%s", toJson(result.Errors, true))
	}
}
