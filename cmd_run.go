package main

import (
	"context"
	"fmt"

	"github.com/cristalhq/dbumper/dbump"
)

type configRun struct {
	Path string `default:"./migrations"`
	DB   string `default:"UNKNOWN"`
}

func runCmd(ctx context.Context) error {
	var cfg configRun

	// TODO: init *sql.DB
	var migrator dbump.Migrator
	switch cfg.DB {
	case "postgres":
		migrator = dbump.NewMigratorPostgres(nil)
	case "mysql":
		migrator = dbump.NewMigratorMySQL(nil)
	default:
		return fmt.Errorf("unsupported DB: %s", cfg.DB)
	}

	return dbump.Run(ctx, migrator, nil)
}
