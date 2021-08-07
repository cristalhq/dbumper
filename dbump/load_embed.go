package dbump

import (
	"embed"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EmbedLoader struct {
	fs   embed.FS
	path string
}

func NewEmbedLoader(fs embed.FS, path string) *EmbedLoader {
	return &EmbedLoader{
		fs:   fs,
		path: path,
	}
}

func (efs *EmbedLoader) Load() ([]*Migration, error) {
	files, err := efs.fs.ReadDir(efs.path)
	if err != nil {
		return nil, err
	}

	migs := make([]*Migration, 0, len(files))
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}

		matches := migrationRE.FindStringSubmatch(fi.Name())
		if len(matches) != 2 {
			continue
		}

		n, err := strconv.ParseInt(matches[1], 10, 32)
		if err != nil {
			return nil, err
		}

		id := int(n)
		switch {
		case id < len(files)+1:
			return nil, fmt.Errorf("duplicate migration %d", id)
		case len(files)+1 < id:
			return nil, fmt.Errorf("missing migration %d", len(files)+1)
		}

		body, err := os.ReadFile(efs.path + fi.Name())
		if err != nil {
			return nil, err
		}

		parts := strings.SplitN(string(body), MigrationDelimiter, 2)
		applySQL := strings.TrimSpace(parts[0])

		var rollbackSQL string
		if len(parts) == 2 {
			rollbackSQL = strings.TrimSpace(parts[1])
		}

		migs = append(migs, &Migration{
			ID:       id,
			Apply:    applySQL,
			Rollback: rollbackSQL,
		})
	}
	return migs, nil
}
