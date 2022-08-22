package main

import (
	"fmt"
	"os"

	"github.com/cristalhq/acmd"
	"github.com/cristalhq/aconfig"
)

var Version = "(devel)"

func main() {
	r := acmd.RunnerOf(cmds, acmd.Config{
		Version: Version,
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}

var cmds = []acmd.Command{
	{
		Name:        "init",
		Description: "initialize migration folder",
		Do:          cmdInitFolder,
	},
	{
		Name:        "new",
		Description: "create a new empty migration",
		Do:          cmdNewMigration,
	},
	{
		Name:        "status",
		Description: "show database status",
		Do:          cmdStatus,
	},
	{
		Name:        "run",
		Description: "run migrations on database",
		Do:          cmdRun,
	},
}

func loadConfig(cfg interface{}) error {
	loader := aconfig.LoaderFor(cfg, aconfig.Config{
		SkipFiles:       true,
		EnvPrefix:       "DBUMPER",
		AllowDuplicates: true,
		Args:            os.Args[2:], // Hack to not propagate os.Args to all commands
	})

	if err := loader.Load(); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	return nil
}
