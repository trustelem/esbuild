package zipfs

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/trustelem/esbuild/pkg/api"
)

type ZipEntry struct {
	Name string
	K    string // "file" or "dir" or "symlink"
	Ctnt []byte
	F    *zip.File
	D    []string
	Link string
}

type ZipFS struct {
	cwd     string
	rootDir string
	root    string
	entries map[string]*ZipEntry
	log     *log.Logger
}

func (*ZipFS) cleanupName(name string) string {
	name = filepath.Clean(name)
	if strings.HasPrefix(name, "./") {
		name = name[2:]
	}
	name = strings.TrimLeft(name, "/")
	name = strings.TrimRight(name, "/")

	return name
}

func (*ZipFS) isRooted(name, root string) bool {
	return name == root ||
		(strings.HasPrefix(name, root) && name[len(root)] == '/')
}

func (*ZipFS) stripRoot(name, root string) string {
	name = strings.TrimPrefix(name, root)
	name = strings.TrimLeft(name, "/")
	return name
}

func New(file string, fsRoot string, zipRoot string, cwd string) (*ZipFS, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, err
	}

	res := &ZipFS{
		rootDir: filepath.Dir(fsRoot),
		root:    fsRoot,
		entries: make(map[string]*ZipEntry),
		cwd:     cwd,
		log:     log.New(io.Discard, "", 0),
	}
	zipRoot = res.cleanupName(zipRoot)

	for _, zf := range zr.File {
		kind := "file"
		mode := zf.Mode() & fs.ModeType
		switch mode {
		case 0:
			kind = "fil"
		case fs.ModeDir:
			kind = "dir"
		case fs.ModeSymlink:
			kind = "sym"

		default:
			panic(fmt.Sprintf("unhandled file mode in zip archive : %d", mode))
		}

		name := res.cleanupName(zf.Name)
		if !res.isRooted(name, zipRoot) {
			continue
		}
		name = res.stripRoot(name, zipRoot)

		entry := res.entries[name]
		if entry != nil && entry.K != kind {
			panic("duplicate entry: " + name + " " + entry.K + " " + kind)
		}
		if entry == nil {
			entry = &ZipEntry{}
			entry.Name = name
			entry.K = kind
			entry.F = zf

			res.entries[name] = entry
		}

		if name != "" {
			dirName := filepath.Dir(name)
			if dirName == "." {
				dirName = ""
			}
			d := res.entries[dirName]
			if d != nil {
				base := filepath.Base(name)
				d.D = append(d.D, base)
			} else {
				log.Printf("no entry found for dir '%s'", dirName)
			}
		}
	}

	return res, nil
}

var _ api.FsLike = &ZipFS{}

func (zfs *ZipFS) Getwd() (string, error) {
	zfs.log.Printf(" ---- zfs.Getwd")
	return zfs.cwd, nil
}

func (zfs *ZipFS) Open(name string) (io.ReadCloser, error) {
	zfs.log.Printf(" ---- zfs.Open %s", name)
	if !strings.HasPrefix(name, zfs.root) {
		return os.Open(name)
	}

	name = zfs.stripRoot(name, zfs.root)

	e := zfs.entries[name]
	if e == nil {
		return nil, fs.ErrNotExist
	}
	if e.K != "fil" {
		return nil, fs.ErrInvalid
	}

	f, err := e.F.Open()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func osReaddirnames(name string) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file.Readdirnames(0)
}

func (zfs *ZipFS) Readdirnames(name string) ([]string, error) {
	zfs.log.Printf(" ---- zfs.Readdirnames %s", name)
	if !strings.HasPrefix(name, zfs.root) {
		names, err := osReaddirnames(name)
		if err == nil && name == zfs.rootDir {
			names = append(names, "node_modules")
		}
		return names, err
	}

	name = zfs.stripRoot(name, zfs.root)

	e := zfs.entries[name]
	if e == nil {
		return nil, fs.ErrNotExist
	}
	if e.K != "dir" {
		return nil, fs.ErrInvalid
	}

	return e.D, nil
}

func (zfs *ZipFS) Stat(path string) (os.FileInfo, error) {
	zfs.log.Printf(" ---- zfs.Stat %s", path)
	if !strings.HasPrefix(path, zfs.root) {
		return os.Stat(path)
	}

	path = zfs.stripRoot(path, zfs.root)

	e := zfs.entries[path]
	if e == nil {
		return nil, fs.ErrNotExist
	}

	return e.F.FileInfo(), nil
}

func (zfs *ZipFS) MkdirAll(path string, perm os.FileMode) error {
	zfs.log.Printf(" ---- zfs.MkdirAll %s", path)
	return os.MkdirAll(path, perm)
}

func (zfs *ZipFS) WriteFile(path string, content []byte, perm os.FileMode) error {
	zfs.log.Printf(" ---- zfs.WriteFile %s", path)
	return os.WriteFile(path, content, perm)
}
