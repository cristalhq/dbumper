package dbump

import (
	"context"
	"database/sql"
)

type MigratorPostgres struct {
	db           *sql.DB
	versionTable string
}

func NewMigratorPostgres(db *sql.DB) *MigratorPostgres {
	return &MigratorPostgres{
		db:           db,
		versionTable: "_schema_version",
	}
}

func (pg *MigratorPostgres) Lock(ctx context.Context) error {
	_, err := pg.db.ExecContext(ctx, "SELECT pg_advisory_lock($1)", lockNum)
	return err
}

func (pg *MigratorPostgres) Unlock(ctx context.Context) error {
	_, err := pg.db.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", lockNum)
	return err
}

func (pg *MigratorPostgres) Version(ctx context.Context) (version int, err error) {
	row := pg.db.QueryRowContext(ctx, "SELECT version FROM "+pg.versionTable)
	err = row.Scan(&version)
	return version, err
}

func (pg *MigratorPostgres) SetVersion(ctx context.Context, version int) error {
	_, err := pg.db.ExecContext(ctx, "UPDATE "+pg.versionTable+" SET version = $1", version)
	return err
}

func (pg *MigratorPostgres) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := pg.db.ExecContext(ctx, query)
	return err
}
