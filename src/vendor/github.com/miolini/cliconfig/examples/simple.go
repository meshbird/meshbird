package main

import (
	"log"
	"os"

	"github.com/miolini/cliconfig"
	"github.com/urfave/cli"
)

type Config struct {
	ListenAddr string `flag:"listen_addr_flag" env:"LISTEN_ADDR_ENV" default:"localhost:8080"`
}

func main() {
	config := Config{}
	app := cli.NewApp()
	app.Name = "example-app"
	app.Flags = cliconfig.Fill(&config, "EXAMPLE_APP_")
	app.Action = func(ctx *cli.Context) error {
		log.Printf("config: %#v", config)
		return nil
	}
	app.Run(os.Args)
}
