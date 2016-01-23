package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gophergala2016/meshbird/common"
	"os/signal"
)

var (
// VERSION var using for auto versioning through Go linker
	VERSION = "dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "meshbird"
	app.Usage = "distributed overlay private networking"
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "create new network",
			Action:  actionNew,
		},
		{
			Name:    "join",
			Aliases: []string{"j"},
			Usage:   "join network",
			Action:  actionJoin,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		println("error: ", err)
	}
}

func actionNew(ctx *cli.Context) {

}

func actionJoin(ctx *cli.Context) {
	node := common.NewLocalNode(nil)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	defer signal.Stop(signalChan)

	err := node.Start()
	if err != nil {
		log.Printf("node start error: %s", err)
	}

	s := <-signalChan
	log.Printf("received signal %s, stopping...", s)
	node.Stop()
}
