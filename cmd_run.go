package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cristalhq/dbumper/dbump"

	_ "github.com/ClickHouse/clickhouse-go" // import ClickHouse
	_ "github.com/go-sql-driver/mysql"      // import MySQL
	_ "github.com/jackc/pgx/v4/stdlib"      // import Postgres (pgx-stdlib)
)

type configRun struct {
	Path string `default:"./migrations"`
	DB   string `default:"UNKNOWN"`
	DSN  string `default:"UNKNOWN"`
}

func runCmd(ctx context.Context) error {
	var cfg configRun
	if err := loadConfig(&cfg); err != nil {
		return err
	}

	var migrator dbump.Migrator
	switch cfg.DB {
	case "clickhouse":
		db, err := sql.Open("clickhouse", cfg.DSN)
		if err != nil {
			return err
		}
		migrator = dbump.NewMigratorClickHouse(db)

	case "mysql":
		db, err := sql.Open("mysql", cfg.DSN)
		if err != nil {
			return err
		}
		migrator = dbump.NewMigratorMySQL(db)

	case "postgres":
		db, err := sql.Open("pgx", cfg.DSN)
		if err != nil {
			return err
		}
		migrator = dbump.NewMigratorPostgres(db)

	default:
		return fmt.Errorf("unsupported DB: %s", cfg.DB)
	}

	return dbump.Run(ctx, migrator, dbump.NewDiskLoader(cfg.Path))
}
