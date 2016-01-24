package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gophergala2016/meshbird/common"
	"github.com/gophergala2016/meshbird/secure"
	"net"
	"os/signal"
	"time"
)

const (
	MeshbirdKeyEnv = "MESHBIRD_KEY"
)

var (
	// VERSION var using for auto versioning through Go linker
	VERSION = "dev"
	logger  = log.New(os.Stderr, "[main] ", log.LstdFlags)
)

func main() {
	go http.ListenAndServe("0.0.0.0:15080", nil)

	app := cli.NewApp()
	app.Name = "meshbird"
	app.Usage = "distributed overlay private networking"
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name:      "new",
			Aliases:   []string{"n"},
			Usage:     "create new network",
			Action:    actionNew,
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
		logger.Printf("error: %s", err)
	}
}

func actionNew(ctx *cli.Context) {
	var secret *secure.NetworkSecret
	var err error

	if len(ctx.Args()) > 0 {
		keyStr := ctx.Args().First()
		secret, err = secure.NetworkSecretUnmarshal(keyStr)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		_, ipnet, err := net.ParseCIDR(ctx.String("CIDR"))
		if err != nil {
			logger.Fatalf("cidr parse error: %s", err)
		}
		secret = secure.NewNetworkSecret(ipnet)
	}
	keyStr := secret.Marshal()
	logger.Printf("key: %s", keyStr)
}

func actionJoin(ctx *cli.Context) {
	key := os.Getenv(MeshbirdKeyEnv)
	if key == "" {
		logger.Fatalf("environment variable %s is not specified", MeshbirdKeyEnv)
	}

	nodeConfig := &common.Config{
		SecretKey: key,
	}
	node, err := common.NewLocalNode(nodeConfig)
	if err != nil {
		logger.Fatalf("local node init error: %s", err)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	defer signal.Stop(signalChan)

	go func() {
		s := <-signalChan
		logger.Printf("received signal %s, stopping...", s)
		node.Stop()

		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	err = node.Start()
	if err != nil {
		logger.Fatalf("node start error: %s", err)
	}

	node.WaitStop()
}
