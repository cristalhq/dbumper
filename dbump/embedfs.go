package dbump

import (
	"embed"
	"os"
)

var _ FS = &EmbedFS{}

type EmbedFS struct {
	fs      embed.FS
	dirname string
}

func NewEmbedFS(fs embed.FS, dirname string) *EmbedFS {
	return &EmbedFS{
		fs:      fs,
		dirname: dirname,
	}
}

func (efs *EmbedFS) ReadDir() ([]os.FileInfo, error) {
	entries, err := efs.fs.ReadDir(efs.dirname)
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

func (efs *EmbedFS) ReadFile(filename string) ([]byte, error) {
	return efs.fs.ReadFile(efs.dirname + filename)
}
