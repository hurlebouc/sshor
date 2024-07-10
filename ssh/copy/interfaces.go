package copy

import (
	"io"
	"io/fs"
)

type writer interface {
	Close() error
	ReadFrom(io.Reader) (int64, error)
}

type reader interface {
	io.Reader
	Close() error
}

type fileSystem interface {
	mkdirAll(path string) error
	create(path string) (writer, error)
	open(path string) (reader, error)
	readDir(path string) ([]fs.FileInfo, error)
	join(a, b string) string
}

type Endpoint struct {
	path       string
	fileSystem fileSystem
}
