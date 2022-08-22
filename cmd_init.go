package main

import (
	"context"
	"fmt"
	"os"
)

type configInit struct {
	Path string `default:"./migrations"`
}

func cmdInitFolder(_ context.Context, _ []string) error {
	var cfg configInit
	if err := loadConfig(&cfg); err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.Path, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}
	return nil
}
