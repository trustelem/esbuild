package fs

import (
	"io"
	"os"
)

type FsLike interface {
	Getwd() (string, error)
	Open(name string) (io.ReadCloser, error)
	Readdirnames(name string) ([]string, error)
	Stat(path string) (os.FileInfo, error)
}
