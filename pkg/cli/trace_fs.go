package cli

import (
	"bytes"
	"io"
	"os"

	"github.com/trustelem/esbuild/internal/fs"
	"github.com/trustelem/esbuild/pkg/api"
)

type TraceFs struct {
	failedirs []string
	dirnames  []string
	openfiles []string

	real fs.FS
}

var _ api.FsLike = &TraceFs{}

func NewTraceFS() (*TraceFs, error) {
	real, err := fs.RealFS(fs.RealFSOptions{})
	if err != nil {
		return nil, err
	}
	return &TraceFs{real: real}, nil
}

func addStr(strs []string, str string) []string {
	for _, s := range strs {
		if s == str {
			return strs
		}
	}

	strs = append(strs, str)
	return strs
}

func (fsys *TraceFs) Getwd() (string, error) {
	return fsys.real.Cwd(), nil
}

func (fsys *TraceFs) Open(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	fsys.openfiles = addStr(fsys.openfiles, name)
	return f, nil
}

func (fsys *TraceFs) Readdirnames(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		fsys.failedirs = addStr(fsys.failedirs, name)
		return nil, err
	}

	fsys.dirnames = addStr(fsys.dirnames, name)
	return f.Readdirnames(0)
}

func (fsys *TraceFs) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (fsys *TraceFs) ScanLog() string {
	if len(fsys.dirnames) == 0 && len(fsys.openfiles) == 0 {
		return "-- nothing was opened --"
	}

	var b bytes.Buffer
	if len(fsys.dirnames) > 0 {
		for _, dir := range fsys.dirnames {
			b.WriteString("\tdir: " + dir + "\n")
		}
		b.WriteString("\n")
	}

	if len(fsys.failedirs) > 0 {
		for _, dir := range fsys.failedirs {
			b.WriteString("\ttry: " + dir + "\n")
		}
		b.WriteString("\n")
	}

	if len(fsys.openfiles) > 0 {
		for _, file := range fsys.openfiles {
			b.WriteString("\tfile: " + file + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
