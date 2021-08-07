package dbump

import (
	"context"
	"database/sql"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// to prevent multiple migrations running at the same time
const lockNum int64 = 777_777_777
