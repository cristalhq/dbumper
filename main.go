package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cristalhq/acmd"
	"github.com/cristalhq/aconfig"
)

var Version = "v0.0.0"

func main() {
	r := acmd.RunnerOf(cmds, acmd.Config{
		Version: Version,
	})

	if err := r.Run(); err != nil {
		log.Fatal(fmt.Errorf("dbumper: %w", err))
	}
}

var cmds = []acmd.Command{
	{
		Name:        "init",
		Description: "initialize migration folder",
		Do:          initFolderCmd,
	},
	{
		Name:        "new",
		Description: "create a new empty migration",
		Do:          newMigrationCmd,
	},
	{
		Name:        "run",
		Description: "run migrations on database",
		Do:          runCmd,
	},
}

var acfg = aconfig.Config{
	SkipFiles:       true,
	EnvPrefix:       "DBUMPER",
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
