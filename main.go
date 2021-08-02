package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/cristalhq/aconfig"
)

func main() {
	ctx, cancel := appContext()
	defer cancel()

	if err := runMain(ctx, os.Args[1:]); err != nil {
		log.Fatal(fmt.Errorf("dbumper: %w", err))
	}
}

type Config struct {
	Init struct {
		Path         string `default:"./migrations"`
		VersionTable string `default:"_dbumper_version"`
		RemoveOld    bool   `default:"false"`
	} `flag:"-"`
	New struct {
		Path string `default:"./migrations"`
		Name string `default:"NONAME"`
	} `flag:"-"`
}

func loadConfig(args []string) (*Config, error) {
	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFiles:       true,
		EnvPrefix:       "dbumper",
		AllowDuplicates: true,
		Args:            args,
	})
	if err := loader.Load(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func runMain(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("no command provided (init, new, run, snapshot, help)")
	}

	cfg, err := loadConfig(args[1:])
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	switch cmd := args[0]; cmd {
	case "init":
		return initFolderCmd(ctx, cfg)
	case "new":
		return newMigrationCmd(ctx, cfg)
	case "run":
		panic("unimplemented")
	case "snapshot":
		panic("unimplemented")
	case "help":
		panic("unimplemented")
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func initFolderCmd(_ context.Context, cfg *Config) error {
	if err := os.Mkdir(cfg.Init.Path, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}
	return nil
}

func newMigrationCmd(_ context.Context, cfg *Config) error {
	path := cfg.New.Path

	migrations, err := loadMigrations(path)
	if err != nil {
		return fmt.Errorf("cannot load migrations: %w", err)
	}

	name := newMigrationFileName(len(migrations)+1, cfg.New.Name)

	if err := createNewMigration(path, name); err != nil {
		return fmt.Errorf("cannot create new migration: %w", err)
	}
	return nil
}

func newMigrationFileName(id int, name string) string {
	return fmt.Sprintf("%04d_%s.sql", id, name)
}

func createNewMigration(path, name string) error {
	path = filepath.Join(path, name)

	newMig, errOpen := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if errOpen != nil {
		return errOpen
	}
	defer newMig.Close()

	_, err := newMig.WriteString(newMigrationText)
	return err
}

func loadMigrations(path string) ([]string, error) {
	path = strings.TrimRight(path, string(filepath.Separator))

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files))
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

		switch id := int(n); {
		case id < len(paths)+1:
			return nil, fmt.Errorf("Duplicate migration %d", n)
		case len(paths)+1 < id:
			return nil, fmt.Errorf("Missing migration %d", len(paths)+1)
		default:
			paths = append(paths, filepath.Join(path, fi.Name()))
		}
	}
	return paths, nil
}

var migrationRE = regexp.MustCompile(`^(\d+)_.+\.sql$`)

var newMigrationText = `--- [Type name of the migration here]

--- apply above / rollback below ---

`

// appContext returns context that will be cancelled on specific OS signals.
func appContext() (context.Context, context.CancelFunc) {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP}

	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	return ctx, cancel
}
