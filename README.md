[![GitHub release](https://img.shields.io/github/release/meshbird/meshbird/all.svg?style=flat-square)](https://github.com/meshbird/meshbird/releases)

# Meshbird 2.0

## About

Meshbird is open-source **cloud-native** multi-region multi-cloud decentralized private networking.

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

## Benchmark

DigitalOcean NYC3 eth1 - private network, encryption on

```# iperf3 -P 1 -R -c 10.247.0.1                                          
Connecting to host 10.247.0.1, port 5201
Reverse mode, remote host 10.247.0.1 is sending
[  4] local 10.247.0.2 port 42536 connected to 10.247.0.1 port 5201
[ ID] Interval           Transfer     Bandwidth
[  4]   0.00-1.00   sec   172 MBytes  1.44 Gbits/sec                  
[  4]   1.00-2.00   sec   162 MBytes  1.36 Gbits/sec                  
[  4]   2.00-3.00   sec   189 MBytes  1.58 Gbits/sec                  
[  4]   3.00-4.00   sec   184 MBytes  1.54 Gbits/sec                  
[  4]   4.00-5.00   sec   176 MBytes  1.48 Gbits/sec                  
[  4]   5.00-6.00   sec   207 MBytes  1.74 Gbits/sec                  
[  4]   6.00-7.00   sec   200 MBytes  1.68 Gbits/sec                  
[  4]   7.00-8.00   sec   192 MBytes  1.61 Gbits/sec                  
[  4]   8.00-9.00   sec   201 MBytes  1.69 Gbits/sec                  
[  4]   9.00-10.00  sec   192 MBytes  1.61 Gbits/sec                  
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bandwidth       Retr
[  4]   0.00-10.00  sec  1.84 GBytes  1.58 Gbits/sec    0             sender
[  4]   0.00-10.00  sec  1.83 GBytes  1.58 Gbits/sec                  receiver```
