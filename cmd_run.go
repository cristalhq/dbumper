package main

import (
	"context"
	"fmt"

	"github.com/cristalhq/dbump"
	"github.com/cristalhq/dbump/dbump_pgx"

	_ "github.com/ClickHouse/clickhouse-go" // import ClickHouse
	_ "github.com/go-sql-driver/mysql"      // import MySQL
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // import Postgres (pgx-stdlib)
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

	migrator, err := getMigrator(ctx, cfg)
	if err != nil {
		return err
	}

	config := dbump.Config{
		Migrator: migrator,
		Loader:   dbump.NewDiskLoader(cfg.Path),
		Mode:     parseMode(cfg.Mode),
		Num:      1,
	}
	return dbump.Run(ctx, config)
}

func getMigrator(ctx context.Context, cfg configRun) (dbump.Migrator, error) {
	switch cfg.DB {
	case "postgres":
		conn, err := pgx.Connect(ctx, cfg.DSN)
		if err != nil {
			return nil, err
		}
		return dbump_pgx.NewMigrator(conn, dbump_pgx.Config{}), nil

	default:
		return nil, fmt.Errorf("unsupported database: %s", cfg.DB)
	}
}

func parseMode(mode string) dbump.MigratorMode {
	switch mode {
	case "up":
		return dbump.ModeApplyAll
	case "down":
		return dbump.ModeRevertAll
	case "up1":
		return dbump.ModeApplyN
	case "down1":
		return dbump.ModeRevertN
	default:
		return dbump.ModeNotSet

	}
}
