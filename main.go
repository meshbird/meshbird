package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/meshbird/meshbird/common"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/secure"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	Version    = "dev"
	NetworkKey string
	LogLevel   string

	keyNotSetError = errors.New("please, set environment variable \"MESHBIRD_KEY\" or specify a flag \"key\"")
)

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
			Name:        "key",
			Usage:       "set network key",
			Destination: &NetworkKey,
			EnvVar:      "MESHBIRD_KEY",
		},
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

func actionGetIP(ctx *cli.Context) error {
	if NetworkKey == "" {
		log.Fatal(keyNotSetError.Error())
	}

	secret, err := secure.NetworkSecretUnmarshal(NetworkKey)
	if err != nil {
		log.Fatal("error on decode network key, %v", err)
	}

	state := common.NewState(secret)
	state.Save()

	fmt.Println(state.PrivateIP.String())
	log.Info("private IP %q restored successfully", state.PrivateIP.String())
	return err
}

func actionNew(ctx *cli.Context) error {
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
	return err
}

func actionJoin(ctx *cli.Context) error {
	if NetworkKey == "" {
		log.Fatal(keyNotSetError.Error())
	}

	node, err := common.NewLocalNode(&common.Config{SecretKey: NetworkKey})
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
	return err
}
