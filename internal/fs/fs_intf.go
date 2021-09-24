package fs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type IntfFS struct {
	cwd   string
	inner FsLike
}

var _ FS = &IntfFS{}

func NewIntfFS(fs FsLike) *IntfFS {
	cwd, _ := fs.Getwd()
	return &IntfFS{
		cwd:   cwd,
		inner: fs,
	}
}

func (f *IntfFS) ReadDirectory(dir string) (entries DirEntries, canonicalError error, originalError error) {
	if !f.IsAbs(dir) {
		dir = f.Join(f.cwd, dir)
	}

	names, err := f.inner.Readdirnames(dir)
	if err != nil {
		return DirEntries{}, err, err
	}

	entries.dir = dir
	entries.data = make(map[string]*Entry)
	for _, name := range names {
		entries.data[strings.ToLower(name)] = &Entry{
			dir:      dir,
			base:     name,
			needStat: true,
		}
	}

	return entries, nil, nil
}

func (f *IntfFS) ReadFile(path string) (contents string, canonicalError error, originalError error) {
	path = f.toAbs(path)

	rc, err := f.inner.Open(path)
	if err != nil {
		return "", err, err
	}

	defer rc.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, rc)
	if err != nil {
		return "", err, err
	}

	return buf.String(), nil, nil
}

func (f *IntfFS) toAbs(foo string) string {
	if !f.IsAbs(foo) {
		foo = f.Join(f.cwd, foo)
	}
	return foo
}

func (f *IntfFS) ModKey(path string) (ModKey, error) {
	path = f.toAbs(path)
	nfo, err := f.inner.Stat(path)
	if err != nil {
		return ModKey{}, err
	}

	size := nfo.Size()
	mode := nfo.Mode()
	ts := nfo.ModTime().Unix()
	tsnano := nfo.ModTime().UnixNano()

	return ModKey{
		size:       size,
		mode:       uint32(mode),
		mtime_sec:  ts,
		mtime_nsec: tsnano,
	}, nil
}

func (f *IntfFS) IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func (f *IntfFS) Abs(path string) (string, bool) {
	abs, err := filepath.Abs(path)
	return abs, err == nil
}

func (f *IntfFS) Dir(path string) string {
	return filepath.Dir(path)
}

func (f *IntfFS) Base(path string) string {
	return filepath.Base(path)
}

func (f *IntfFS) Ext(path string) string {
	return filepath.Ext(path)
}

func (f *IntfFS) Join(parts ...string) string {
	return filepath.Join(parts...)
}

func (f *IntfFS) Cwd() string {
	wd, _ := os.Getwd()
	return wd
}

func (f *IntfFS) Rel(base string, target string) (string, bool) {
	path, err := filepath.Rel(base, target)
	return path, err == nil
}

func (f *IntfFS) kind(dir string, base string) (symlink string, kind EntryKind) {
	dir = f.toAbs(dir)
	abs := f.Join(dir, base)

	nfo, err := f.inner.Stat(abs)
	if err != nil {
		return "", 0
	}

	if nfo.IsDir() {
		return "", DirEntry
	}
	return "", FileEntry
}

func (f *IntfFS) OpenFile(path string) (result OpenedFile, canonicalError error, originalError error) {
	panic("OpenFile not implemented")
	return nil, nil, nil
}

func (f *IntfFS) WatchData() WatchData {
	//helpers.DumpStack(os.Stderr, "----- WatchData called")
	return WatchData{}
}
