package dbump

import (
	"context"
	"database/sql"
)

type MigratorMySQL struct {
	db           *sql.DB
	versionTable string
}

func NewMigratorMySQL(db *sql.DB) *MigratorMySQL {
	return &MigratorMySQL{
		db:           db,
		versionTable: "_schema_version",
	}
}

func (my *MigratorMySQL) Lock(ctx context.Context) error {
	_, err := my.db.ExecContext(ctx, `SELECT GET_LOCK(?, 10)`, lockNum)
	return err
}

func (my *MigratorMySQL) Unlock(ctx context.Context) error {
	_, err := my.db.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", lockNum)
	return err
}

func (my *MigratorMySQL) Version(ctx context.Context) (version int, err error) {
	row := my.db.QueryRowContext(ctx, "SELECT version FROM "+my.versionTable)
	err = row.Scan(&version)
	return version, err
}

func (my *MigratorMySQL) SetVersion(ctx context.Context, version int) error {
	_, err := my.db.ExecContext(ctx, "UPDATE "+my.versionTable+" SET version = $1", version)
	return err
}

func (my *MigratorMySQL) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := my.db.ExecContext(ctx, query)
	return err
}
