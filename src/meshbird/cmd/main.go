package main

import (
	"log"
	"meshbird"
	"meshbird/config"

	"github.com/alexflint/go-arg"
)

func main() {
	var cfg config.Config
	arg.MustParse(&cfg)
	log.Printf("config: %#v", cfg)
	app := meshbird.NewApp(cfg)
	err := app.Run()
	if err != nil {
		log.Printf("run err: %s", err)
	}
}
