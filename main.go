package main

import (
	"log"

	"github.com/meshbird/meshbird/common"

	"github.com/alexflint/go-arg"
)

func main() {
	var cfg common.Config
	arg.MustParse(&cfg)
	log.Printf("config: %#v", cfg)
	app := common.NewApp(cfg)
	err := app.Run()
	if err != nil {
		log.Printf("run err: %s", err)
	}
}
