package dbump

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var migrationRE = regexp.MustCompile(`^(\d+)_.+\.sql$`)

// DiskLoader can load migrations from disk/OS.
type DiskLoader struct {
	path string
}

// NewDiskLoader instantiates a new DiskLoader.
func NewDiskLoader(path string) *DiskLoader {
	return &DiskLoader{
		path: strings.TrimRight(path, string(os.PathSeparator)),
	}
}

// Load is a method for Loader interface.
func (fs *DiskLoader) Load() ([]*Migration, error) {
	files, err := os.ReadDir(fs.path)
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
		case id < len(migs)+1:
			return nil, fmt.Errorf("duplicate migration %d", id)
		case len(migs)+1 < id:
			return nil, fmt.Errorf("missing migration %d", len(migs)+1)
		}

		body, err := os.ReadFile(filepath.Join(fs.path, fi.Name()))
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
			ID: id,
			// Name:     filename,
			Apply:    applySQL,
			Rollback: rollbackSQL,
		})
	}
	return migs, nil
}
