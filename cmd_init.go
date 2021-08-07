package main

import (
	"context"
	"fmt"
	"os"
)

type ConfigInit struct {
	Path         string `default:"./migrations"`
	VersionTable string `default:"_dbumper_version"`
	RemoveOld    bool   `default:"false"`
}

func initFolderCmd(_ context.Context) error {
	var cfg ConfigInit
	if err := loadConfig(&cfg); err != nil {
		return err
	}

	if cfg.RemoveOld {
		if err := os.RemoveAll(cfg.Path); err != nil {
			return fmt.Errorf("cannot remove directory: %w", err)
		}
	}
	if err := os.Mkdir(cfg.Path, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}
	return nil
}