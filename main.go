package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/meshbird/meshbird/common"
	"github.com/meshbird/meshbird/secure"
	"net"
	"os"
	"os/signal"
	"time"
)

const (
	MeshbirdKeyEnv = "MESHBIRD_KEY"
)

var (
	// VERSION var using for auto versioning through Go linker
	VERSION  = "dev"
	logger   = log.New()
	Loglevel = ""
)

func main() {
	app := cli.NewApp()
	app.Name = "MeshBird"
	app.Usage = "distributed private networking"
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "loglevel",
			Usage:       "log level",
			Destination: &Loglevel,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithError(err).Fatal()
	}
}

func actionGetIP(ctx *cli.Context) {
	keyStr := os.Getenv(MeshbirdKeyEnv)
	if keyStr == "" {
		logger.Fatalf("environment variable %s is not specified", MeshbirdKeyEnv)
	}
	secret, err := secure.NetworkSecretUnmarshal(keyStr)
	if err != nil {
		logger.WithError(err).Fatal()
	}
	state := common.NewState(secret)
	state.Save()
	logger.Info("Restored private IP address %s from state", state.PrivateIP.String())
}

func actionNew(ctx *cli.Context) {
	var secret *secure.NetworkSecret
	var err error

	if len(ctx.Args()) > 0 {
		keyStr := ctx.Args().First()
		secret, err = secure.NetworkSecretUnmarshal(keyStr)
		if err != nil {
			logger.WithError(err).Fatal()
		}
	} else {
		_, ipnet, err := net.ParseCIDR(ctx.String("CIDR"))
		if err != nil {
			logger.WithError(err).Fatal("cidr parse error")
		}
		secret = secure.NewNetworkSecret(ipnet)
	}
	fmt.Println(secret.Marshal())
}

func actionJoin(ctx *cli.Context) {
	logLevel, err := log.ParseLevel(Loglevel)
	if err != nil {
		logger.WithError(err).Fatal()
	}
	logger.Level = logLevel

	key := os.Getenv(MeshbirdKeyEnv)
	if key == "" {
		logger.Fatalf("environment variable %s is not specified", MeshbirdKeyEnv)
	}

	nodeConfig := &common.Config{
		SecretKey: key,
		Loglevel:  logLevel,
	}
	node, err := common.NewLocalNode(nodeConfig)
	if err != nil {
		logger.WithError(err).Fatal("local node init error")
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	defer signal.Stop(signalChan)

	go func() {
		s := <-signalChan
		logger.WithField("signal", s).Info("signal received, stopping...")
		node.Stop()

		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	err = node.Start()
	if err != nil {
		logger.WithError(err).Fatal("node start error")
	}

	node.WaitStop()
}
