package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cristalhq/dbump"

	_ "github.com/ClickHouse/clickhouse-go" // import ClickHouse
	_ "github.com/go-sql-driver/mysql"      // import MySQL
	_ "github.com/jackc/pgx/v4/stdlib"      // import Postgres (pgx-stdlib)
)

type configRun struct {
	Path string `default:"./migrations"`
	DB   string `default:"UNKNOWN"`
	DSN  string `default:"UNKNOWN"`
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
	return dbump.Run(ctx, migrator, dbump.NewDiskLoader(cfg.Path))
}

func getMigrator(cfg configRun) (dbump.Migrator, error) {
	switch cfg.DB {
	case "clickhouse":
		db, err := sql.Open("clickhouse", cfg.DSN)
		if err != nil {
			return nil, err
		}
		return dbump.NewMigratorClickHouse(db), nil

	case "mysql":
		db, err := sql.Open("mysql", cfg.DSN)
		if err != nil {
			return nil, err
		}
		return dbump.NewMigratorMySQL(db), nil

	case "postgres":
		db, err := sql.Open("pgx", cfg.DSN)
		if err != nil {
			return nil, err
		}
		return dbump.NewMigratorPostgres(db), nil

	default:
		return nil, fmt.Errorf("unsupported DB: %s", cfg.DB)
	}
}
