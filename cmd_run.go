package main

import (
	"context"
	"fmt"

	"github.com/cristalhq/dbumper/dbump"
	"github.com/cristalhq/dbumper/dbumpmysql"
	"github.com/cristalhq/dbumper/dbumppg"
)

type ConfigRun struct {
	Path string `default:"./migrations"`
	DB   string `default:"UNKNOWN"`
}

func runCmd(ctx context.Context) error {
	var cfg *ConfigRun

	// TODO: init *sql.DB
	var migrator dbump.Migrator
	switch cfg.DB {
	case "postgres":
		migrator = dbumppg.NewMigrator(nil)
	case "mysql":
		migrator = dbumpmysql.NewMigrator(nil)
	default:
		return fmt.Errorf("unsupported DB: %s", cfg.DB)
	}

	return dbump.Run(ctx, migrator, dbump.NewRealFS(cfg.Path))
}
