[![GitHub release](https://img.shields.io/github/release/meshbird/meshbird/all.svg?style=flat-square)](https://github.com/meshbird/meshbird/releases)

# Meshbird 2.0

## About

Multi-region multi-cloud decentralized private networking.

## Install

Download and install latest release from this page https://github.com/meshbird/meshbird/releases.

## How to run

```
$ ./meshbird --help
Password:
NAME:
   meshbird - multi-region multi-cloud decentralized private networking. 

USAGE:
   meshbird [global options] command [command options] [arguments...]

VERSION:
   2.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --key value                 (default: "hello-world") [$MESHBIRD_KEY]
   --seed_addrs value          (default: "dc1/10.0.0.1/16,dc2/10.0.0.2") [$MESHBIRD_SEED_ADDRS]
   --local_addr value          (default: "10.0.0.1") [$MESHBIRD_LOCAL_ADDR]
   --local_private_addr value  (default: "192.168.0.1") [$MESHBIRD_LOCAL_PRIVATE_ADDR]
   --dc value                  (default: "dc1") [$MESHBIRD_DC]
   --transport_threads value   (default: 1) [$MESHBIRD_TRANSPORT_THREADS]
   --ip value                  (default: "10.237.0.1/16") [$MESHBIRD_IP]
   --mtu value                 (default: 9000) [$MESHBIRD_MTU]
   --verbose value             (default: 0) [$MESHBIRD_VERBOSE]
   --help, -h                  show help
   --version, -v               print the version
```