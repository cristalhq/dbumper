package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cristalhq/dbump"
)

type configNew struct {
	Path string `default:"./migrations"`
	Name string `default:"NONAME"`
}

var newMigrationText = `` +
	`--- [Type name of the migration here]` +
	"\n\n" +
	dbump.MigrationDelimiter +
	"\n"

func newMigrationCmd(_ context.Context, _ []string) error {
	var cfg configNew
	if err := loadConfig(&cfg); err != nil {
		return err
	}
	path := cfg.Path

	migrations, err := dbump.NewDiskLoader(path).Load()
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
