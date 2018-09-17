cliconfig
=========

Urfave Cli flags setup from config struct fields

[![Build Status](https://travis-ci.org/miolini/cliconfig.svg?branch=master)](https://travis-ci.org/miolini/cliconfig)
[![GoDoc](https://godoc.org/github.com/miolini/cliconfig?status.svg)](https://godoc.org/github.com/miolini/cliconfig)
[![Go Report Card](https://goreportcard.com/badge/miolini/cliconfig)](https://goreportcard.com/report/miolini/cliconfig)

## Example

```go
package main

import (
  "log"
  "os"

  "github.com/miolini/cliconfig"
  "github.com/urfave/cli"
)

type Config struct {
  ListenAddr string `flag:"listen_addr_flag" env:"LISTEN_ADDR_ENV" default:"localhost:8080"`
}

func main() {
  config := Config{}
  app := cli.NewApp()
  app.Name = "example-app"
  app.Flags = cliconfig.Fill(&config, "EXAMPLE_APP_")
  app.Action = func(ctx *cli.Context) error {
    log.Printf("config: %#v", config)
    return nil
  }
  app.Run(os.Args)
}
```

Usage

```shell
$ ./simple  --help
command-line-arguments
NAME:
   example-app - A new cli application

USAGE:
   simple [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --listen_addr_flag value  (default: "localhost:8080") [$EXAMPLE_APP_LISTEN_ADDR_ENV]
   --help, -h                show help
   --version, -v             print the version
```

