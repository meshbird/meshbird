package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/meshbird/meshbird/common"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/secure"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	Version    = "dev"
	NetworkKey string
	LogLevel   string

	keyNotSetError = errors.New("please, set environment variable \"MESHBIRD_KEY\"")
)

func init() {
	if envVal := os.Getenv("MESHBIRD_KEY"); envVal != "" {
		NetworkKey = strings.TrimSpace(envVal)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "MeshBird"
	app.Usage = "distributed private networking"
	app.Version = Version
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
					Value: "10.7.0.0/16",
					Usage: "Define custom CIDR",
				},
			},
		},
		{
			Name:    "join",
			Aliases: []string{"j"},
			Usage:   "join network",
			Action:  actionJoin,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bootstrap",
					Value: "",
					Usage: "Define bootstrap nodes for DHT",
				},
			},
		},
		{
			Name:      "ip",
			Aliases:   []string{"i"},
			Usage:     "init state",
			Action:    actionGetIP,
			ArgsUsage: "<key>",
		},
	}
	app.Before = func(context *cli.Context) error {
		log.SetLevel(log.MustParseLevel(LogLevel))
		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "loglevel",
			Value:       "warning",
			Usage:       "set log level",
			Destination: &LogLevel,
			EnvVar:      "MESHBIRD_LOG_LEVEL",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("error on run app, %v", err)
	}
}

func actionGetIP(ctx *cli.Context) {
	if NetworkKey == "" {
		log.Fatal(keyNotSetError.Error())
	}

	secret, err := secure.NetworkSecretUnmarshal(NetworkKey)
	if err != nil {
		log.Fatal("error on decode network key, %v", err)
	}

	state := common.NewState(secret)
	if err = state.Save(); err != nil {
		log.Fatal("state save err: %s", err)
	}

	fmt.Println(state.PrivateIP().String())
	log.Info("private IP %q restored successfully", state.PrivateIP().String())
}

func actionNew(ctx *cli.Context) {
	var secret *secure.NetworkSecret
	var err error

	if len(ctx.Args()) > 0 {
		keyStr := ctx.Args().First()
		secret, err = secure.NetworkSecretUnmarshal(keyStr)
		if err != nil {
			log.Fatal("error on decode network key, %v", err)
		}
	} else {
		_, ipNet, err := net.ParseCIDR(ctx.String("CIDR"))
		if err != nil {
			log.Fatal("cidr parse error, %v", err)
		}
		secret = secure.NewNetworkSecret(ipNet)
	}
	fmt.Println(secret.Marshal())
}

func actionJoin(ctx *cli.Context) {
	if NetworkKey == "" {
		log.Fatal(keyNotSetError.Error())
	}

	config := &common.Config{
		SecretKey:      NetworkKey,
		BootstrapNodes: ctx.String("bootstrap"),
	}

	node, err := common.NewLocalNode(config)
	if err != nil {
		log.Fatal("error on setup local node, %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	defer signal.Stop(signalChan)

	go func() {
		<-signalChan
		log.Debug("signal received, stopping...")
		node.Stop()

		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	err = node.Start()
	if err != nil {
		log.Fatal("error on local node start, %v", err)
	}

	node.WaitStop()
}
