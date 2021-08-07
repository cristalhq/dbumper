package dbump

import (
	"context"
	"fmt"
)

// MigrationDelimiter separates apply and rollback queries inside a migration step/file.
const MigrationDelimiter = `--- apply above / rollback below ---`

// Migrator represents DB over which we will run migration queries.
type Migrator interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error

	Version(ctx context.Context) (version int, err error)
	SetVersion(ctx context.Context, version int) error

	Exec(ctx context.Context, query string, args ...interface{}) error
}

// Loader returns migrations to be applied on a DB.
type Loader interface {
	Load() ([]*Migration, error)
}

// Migration represents migration step that will be runned on DB.
type Migration struct {
	ID       int    // ID of the migration, unique, positive, starts from 1.
	Name     string // Name of the migration
	Apply    string // Apply query
	Rollback string // Rollback query
}

// Run the Migrator with migration queries provided by the Loader.
func Run(ctx context.Context, m Migrator, l Loader) error {
	migs, err := l.Load()
	if err != nil {
		return err
	}
	return runMigration(ctx, m, migs)
}

func runMigration(ctx context.Context, m Migrator, migs []*Migration) error {
	if err := m.Lock(ctx); err != nil {
		return err
	}

	var err error
	defer func() {
		if errUnlock := m.Unlock(ctx); err == nil && errUnlock != nil {
			err = errUnlock
		}
	}()

	err = runLockedMigration(ctx, m, migs)
	return err
}

func runLockedMigration(ctx context.Context, m Migrator, migs []*Migration) error {
	currentVersion, err := m.Version(ctx)
	if err != nil {
		return err
	}

	// TODO: configure
	targetVersion := len(migs)
	switch {
	case targetVersion < 0 || len(migs) < targetVersion:
		fallthrough
	case currentVersion < 0 || len(migs) < currentVersion:
		return fmt.Errorf("target version %d is outside of range 0..%d ", targetVersion, len(migs))
	}

	direction := 1
	if currentVersion > targetVersion {
		direction = -1
	}

	for currentVersion != targetVersion {
		current := migs[currentVersion]
		sequence := current.ID
		query := current.Apply

		if direction == -1 {
			current = migs[currentVersion-1]
			sequence = current.ID - 1
			query = current.Rollback
		}

		if err := m.Exec(ctx, query); err != nil {
			return err
		}

		if err := m.SetVersion(ctx, sequence); err != nil {
			return err
		}
		currentVersion += direction
	}
	return nil
}
