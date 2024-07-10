package copy

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type localFS struct{}

func (local localFS) join(a, b string) string {
	return filepath.Join(a, b)
}

func (local localFS) mkdirAll(path string) error {
	return os.MkdirAll(path, os.ModeDir)
}

func (local localFS) create(path string) (writer, error) {
	return os.Create(path)
}

func (local localFS) open(path string) (reader, error) {
	return os.Open(path)
}

func (local localFS) readDir(path string) ([]fs.FileInfo, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	res := make([]fs.FileInfo, len(dirEntries))
	for i, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		res[i] = info
	}
	return res, nil
}

func (local localFS) isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return stat.IsDir()
}

func (local localFS) url(path string) string {
	return path
}

func (local localFS) exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic(err)
	}
}

func (local localFS) close() {}

func NewLocal(path string) Endpoint {
	return Endpoint{
		path:       path,
		fileSystem: localFS{},
	}
}
