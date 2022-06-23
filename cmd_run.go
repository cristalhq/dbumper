package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cristalhq/dbump"
	"github.com/cristalhq/dbump/dbump_pg"

	_ "github.com/ClickHouse/clickhouse-go" // import ClickHouse
	_ "github.com/go-sql-driver/mysql"      // import MySQL
	_ "github.com/jackc/pgx/v4/stdlib"      // import Postgres (pgx-stdlib)
)

type configRun struct {
	Path string `default:"./migrations"`
	DB   string `default:"UNKNOWN"`
	DSN  string `default:"UNKNOWN"`
	Mode string `default:"UNKNOWN"`
}

func runCmd(ctx context.Context, _ []string) error {
	var cfg configRun
	if err := loadConfig(&cfg); err != nil {
		return err
	}

	migrator, err := getMigrator(cfg)
	if err != nil {
		return err
	}

	config := dbump.Config{
		Migrator: migrator,
		Loader:   dbump.NewDiskLoader(cfg.Path),
		Mode:     parseMode(cfg.Mode),
	}
	return dbump.Run(ctx, config)
}

func getMigrator(cfg configRun) (dbump.Migrator, error) {
	switch cfg.DB {
	case "postgres":
		db, err := sql.Open("pgx", cfg.DSN)
		if err != nil {
			return nil, err
		}
		return dbump_pg.NewMigrator(db), nil

	default:
		return nil, fmt.Errorf("unsupported database: %s", cfg.DB)
	}
}

func parseMode(mode string) dbump.MigratorMode {
	switch mode {
	case "Up":
		return dbump.ModeUp
	case "Down":
		return dbump.ModeDown
	case "UpOne":
		return dbump.ModeUpOne
	case "DownOne":
		return dbump.ModeDownOne
	default:
		return dbump.ModeNotSet

	}
}
