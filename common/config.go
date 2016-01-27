package common

import "github.com/Sirupsen/logrus"

type Config struct {
	SecretKey string
	NetworkID string
	Loglevel  logrus.Level
}
