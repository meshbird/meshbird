package main

import (
	"log"
	"os"

	"meshbird"
	"meshbird/config"

	"github.com/miolini/cliconfig"
	"github.com/urfave/cli"
)

func main() {
	var cfg config.Config
	app := cli.NewApp()
	app.Name = "meshbird"
	app.Flags = cliconfig.Fill(&cfg, "MESHBIRD_")
	app.Action = func(ctx *cli.Context) error {
		log.Printf("config: %#v", cfg)
		meshbirdApp := meshbird.NewApp(cfg)
		err := meshbirdApp.Run()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Printf("app run err: %s", err)
	}
}
