package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cristalhq/dbump"
)

type configStatus struct {
	Path string `default:"./migrations"`
	DB   string `default:"UNKNOWN"`
	DSN  string `default:"UNKNOWN"`

	output io.Writer
}

func cmdStatus(ctx context.Context, _ []string) error {
	var cfg configStatus
	if err := loadConfig(&cfg); err != nil {
		return err
	}
	cfg.output = os.Stdin

	fmt.Fprintln(cfg.output, "Loading migrations from:", cfg.Path)

	loader := dbump.NewDiskLoader(cfg.Path)
	migrations, err := loader.Load()
	if err != nil {
		return err
	}

	mig, err := getMigrator(ctx, configRun{DB: cfg.DB, DSN: cfg.DSN})
	if err != nil {
		return err
	}
	version, err := mig.Version(ctx)
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		if version != -1 {
			fmt.Fprintln(cfg.output, "have version but not migrations")
			return errors.New("consistency")
		}
		fmt.Fprintln(cfg.output, "no migrations")
		return nil
	}

	curr := 1
	if version > 0 {
		fmt.Fprintln(cfg.output, "Applied:")
		for ; curr <= version; curr++ {
			fmt.Fprintln(cfg.output, curr, migrations[curr-1].Name)
		}
	}

	if curr <= len(migrations) {
		fmt.Fprintln(cfg.output, "Not applied:")
		for ; curr <= len(migrations); curr++ {
			fmt.Fprintln(cfg.output, curr, migrations[curr-1].Name)
		}
	}
	return nil
}
