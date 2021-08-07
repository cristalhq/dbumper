package dbump

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// MigrationDelimiter ...
const MigrationDelimiter = `--- apply above / rollback below ---`

// Migrator ...
type Migrator interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error

	Version(ctx context.Context) (version int, err error)
	SetVersion(ctx context.Context, version int) error

	Exec(ctx context.Context, query string, args ...interface{}) error
}

// FS ...
type FS interface {
	ReadDir() ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

// Migration ...
type Migration struct {
	ID       int
	File     string
	Apply    string
	Rollback string
}

// Run ...
func Run(ctx context.Context, m Migrator, fs FS) error {
	migs, err := loadMigrations(fs)
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

	var direction int
	if currentVersion < targetVersion {
		direction = 1
	} else {
		direction = -1
	}

	for currentVersion != targetVersion {
		var current *Migration
		var sequence int
		var query string

		if direction == 1 {
			current = migs[currentVersion]
			sequence = current.ID
			query = current.Apply
		} else {
			current = migs[currentVersion-1]
			sequence = current.ID - 1
			query = current.Rollback
			if current.Rollback == "" {
				return fmt.Errorf("no rollback downgrade of %d", sequence)
			}
		}

		if err := m.Exec(ctx, query); err != nil {
			return err
		}

		if err := m.SetVersion(ctx, sequence); err != nil {
			return err
		}
		currentVersion = currentVersion + direction
	}
	return nil
}

var migrationRE = regexp.MustCompile(`^(\d+)_.+\.sql$`)

func loadMigrations(fs FS) ([]*Migration, error) {
	files, err := fs.ReadDir()
	if err != nil {
		return nil, err
	}

	migs := make([]*Migration, 0, len(files))
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}

		matches := migrationRE.FindStringSubmatch(fi.Name())
		if len(matches) != 2 {
			continue
		}

		n, err := strconv.ParseInt(matches[1], 10, 32)
		if err != nil {
			return nil, err
		}

		id := int(n)
		switch {
		case id < len(files)+1:
			return nil, fmt.Errorf("duplicate migration %d", id)
		case len(files)+1 < id:
			return nil, fmt.Errorf("missing migration %d", len(files)+1)
		}

		mig, err := loadMigration(fs, fi.Name())
		if err != nil {
			return nil, err
		}

		mig.ID = id
		migs = append(migs, mig)
	}
	return migs, nil
}

func loadMigration(fs FS, filename string) (*Migration, error) {
	body, err := fs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(string(body), MigrationDelimiter, 2)
	applySQL := strings.TrimSpace(parts[0])

	var rollbackSQL string
	if len(parts) == 2 {
		rollbackSQL = strings.TrimSpace(parts[1])
	}

	return &Migration{
		File:     filename,
		Apply:    applySQL,
		Rollback: rollbackSQL,
	}, nil
}
