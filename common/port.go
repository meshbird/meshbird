package common

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func GetRandomPort(args ...int) int {
	var min, max int = 4000, 60000
	if len(args) > 0 {
		min = args[0]
	}
	if len(args) > 1 {
		max = args[1]
	}
	return rand.Intn(max-min) + min
}
