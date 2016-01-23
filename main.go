package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gophergala2016/meshbird/common"
	"github.com/gophergala2016/meshbird/ecdsa"
	"os/signal"
)

const (
	MeshbirdKeyEnv = "MESHBIRD_KEY"
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
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "CIDR",
					Value: "192.168.137.1/24",
					Usage: "Define custom CIDR",
				},
			},
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
	key := new(ecdsa.Key)
	if len(ctx.Args())>0 {
		key = ecdsa.Unpack([]byte(ctx.Args()[0]))
	} else {
		key,_ = ecdsa.GenerateKey()
		key.CIDR = ctx.String("CIDR")
	}
	println(string(ecdsa.Pack(key)))
}

func actionJoin(ctx *cli.Context) {
	key := os.Getenv(MeshbirdKeyEnv)
	if key == "" {
		log.Fatalf("environment variable %s is not specified", MeshbirdKeyEnv)
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
