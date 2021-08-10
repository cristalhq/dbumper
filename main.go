package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
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

func runMain(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("no command provided (init, new, run, snapshot, help)")
	}

	switch cmd := args[0]; cmd {
	case "init":
		return initFolderCmd(ctx)
	case "new":
		return newMigrationCmd(ctx)
	case "run":
		return runCmd(ctx)
	case "snapshot":
		panic("unimplemented")
	case "version":
		return runVersion(ctx)
	case "help":
		panic("unimplemented")
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}

type _ struct {
	configInit
	configNew
	configRun
}

var acfg = aconfig.Config{
	SkipFiles:       true,
	EnvPrefix:       "dbumper",
	AllowDuplicates: true,
	Args:            os.Args[2:], // Hack to not propagate os.Args to all commands
}

func loadConfig(cfg interface{}) error {
	loader := aconfig.LoaderFor(cfg, acfg)

	if err := loader.Load(); err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}
	return nil
}

// appContext returns context that will be cancelled on specific OS signals.
func appContext() (context.Context, context.CancelFunc) {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP}

	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	return ctx, cancel
}
