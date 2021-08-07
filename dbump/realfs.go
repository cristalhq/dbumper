package dbump

import (
	"os"
	"path/filepath"
	"strings"
)

var _ FS = &RealFS{}

type RealFS struct {
	path string
}

func NewRealFS(path string) *RealFS {
	return &RealFS{
		path: strings.TrimRight(path, string(filepath.Separator)),
	}
}

func (fs *RealFS) ReadDir() ([]os.FileInfo, error) {
	entries, err := os.ReadDir(fs.path)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]os.FileInfo, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			return nil, err
		}
		fileInfos = append(fileInfos, info)
	}
	return fileInfos, nil
}

func (fs *RealFS) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(fs.path + filename)
}
