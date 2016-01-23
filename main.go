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
			ArgsUsage: "<key>",
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
		log.Printf("error: %s", err)
	}
}

func actionNew(ctx *cli.Context) {

}

func actionJoin(ctx *cli.Context) {
	log.Printf("flags: %s", ctx.FlagNames())
	key := os.Getenv("MESHBIRD_KEY")
	if key == "" {
		log.Fatal("key is not specified")
	}

	nodeConfig := &common.Config{
		SecretKey: key,
	}
	node := common.NewLocalNode(nodeConfig)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	defer signal.Stop(signalChan)

	go func() {
		s := <-signalChan
		log.Printf("received signal %s, stopping...", s)
		node.Stop()
	}()

	err := node.Start()
	if err != nil {
		log.Printf("node start error: %s", err)
	}

	node.WaitStop()
}
