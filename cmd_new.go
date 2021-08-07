package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type configNew struct {
	Path string `default:"./migrations"`
	Name string `default:"NONAME"`
}

var migrationRE = regexp.MustCompile(`^(\d+)_.+\.sql$`)

var newMigrationText = `--- [Type name of the migration here]

--- apply above / rollback below ---

`

func newMigrationCmd(_ context.Context) error {
	var cfg configNew
	if err := loadConfig(&cfg); err != nil {
		return err
	}
	path := cfg.Path

	migrations, err := loadMigrations(path)
	if err != nil {
		return fmt.Errorf("cannot load migrations: %w", err)
	}

	name := newMigrationFileName(len(migrations)+1, cfg.Name)

	if err := createNewMigration(path, name); err != nil {
		return fmt.Errorf("cannot create new migration: %w", err)
	}
	return nil
}

func newMigrationFileName(id int, name string) string {
	return fmt.Sprintf("%04d_%s.sql", id, name)
}

func createNewMigration(path, name string) error {
	path = filepath.Join(path, name)

	newMig, errOpen := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if errOpen != nil {
		return errOpen
	}
	defer newMig.Close()

	_, err := newMig.WriteString(newMigrationText)
	return err
}

// TODO: reuse dbump.Loader
func loadMigrations(path string) ([]string, error) {
	path = strings.TrimRight(path, string(filepath.Separator))

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files))
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

		switch id := int(n); {
		case id < len(paths)+1:
			return nil, fmt.Errorf("duplicate migration %d", id)
		case len(paths)+1 < id:
			return nil, fmt.Errorf("missing migration %d", len(paths)+1)
		default:
			paths = append(paths, filepath.Join(path, fi.Name()))
		}
	}
	return paths, nil
}
