package copy

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
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
	isDir(path string) bool
	url(path string) string
	exists(path string) bool
	close()
}

type Endpoint struct {
	path       string
	fileSystem fileSystem
}

func (e Endpoint) join(path string) Endpoint {
	return Endpoint{
		path:       e.fileSystem.join(e.path, path),
		fileSystem: e.fileSystem,
	}
}

func (e Endpoint) Close() {
	e.fileSystem.close()
}

func (e Endpoint) isDir() bool {
	return e.fileSystem.isDir(e.path)
}

func (e Endpoint) isFile() bool {
	return !e.isDir()
}

func (e Endpoint) name() string {
	return filepath.Base(e.path) //todo : revoir si base sait indifféremment fonctionner avec windows et unix.
}

func (e Endpoint) url() string {
	return e.fileSystem.url(e.path)
}
func (e Endpoint) exists() bool {
	return e.fileSystem.exists(e.path)
}

func CompleteDstPath(src, dst Endpoint) Endpoint {
	if !dst.exists() {
		return dst // incohérent avec le reste, mais même logique que cp
	}
	if src.isFile() && dst.isFile() {
		return dst
	}
	if src.isFile() && dst.isDir() {
		return dst.join(src.name())
	}
	if src.isDir() && dst.isFile() {
		panic(fmt.Sprintf("cannot copy directory %s to file %s", src.url(), dst.url()))
	}
	return dst.join(src.name())
}
