package dbumpmysql

import (
	"context"
	"database/sql"

	"github.com/cristalhq/dbumper/dbump"
)

var _ dbump.Migrator = &Migrator{}

type Migrator struct {
	db           *sql.DB
	versionTable string
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db:           db,
		versionTable: "_schema_version",
	}
}

// to prevent multiple migrations running at the same time
const lockNum int64 = 777_777_777

func (m *Migrator) Lock(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, `SELECT GET_LOCK(?, 10)`, lockNum)
	return err
}

func (m *Migrator) Unlock(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", lockNum)
	return err
}

func (m *Migrator) Version(ctx context.Context) (version int, err error) {
	row := m.db.QueryRowContext(ctx, "SELECT version FROM "+m.versionTable)
	err = row.Scan(&version)
	return version, err
}

func (m *Migrator) SetVersion(ctx context.Context, version int) error {
	_, err := m.db.ExecContext(ctx, "UPDATE "+m.versionTable+" SET version = $1", version)
	return err
}

func (m *Migrator) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := m.db.ExecContext(ctx, query)
	return err
}
