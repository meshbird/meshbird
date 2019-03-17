[![GitHub release](https://img.shields.io/github/release/meshbird/meshbird/all.svg?style=flat-square)](https://github.com/meshbird/meshbird/releases)

# Meshbird 2

## About

Meshbird is open-source **cloud-native** multi-region multi-cloud decentralized private networking.

## Install

Download and install latest release from this page https://github.com/meshbird/meshbird/releases.

## How to run

```
$ ./meshbird --help
Usage: meshbird [--key KEY] [--hostaddr HOSTADDR] [--publicaddrs PUBLICADDRS] [--bindaddrs BINDADDRS] [--seedaddrs SEEDADDRS] [--transportthreads TRANSPORTTHREADS] --ip IP [--mtu MTU] [--verbose VERBOSE]

Options:
  --key KEY
  --hostaddr HOSTADDR
  --publicaddrs PUBLICADDRS, -p PUBLICADDRS
  --bindaddrs BINDADDRS, -b BINDADDRS
  --seedaddrs SEEDADDRS, -s SEEDADDRS
  --transportthreads TRANSPORTTHREADS
  --ip IP
  --mtu MTU
  --verbose VERBOSE
  --help, -h             display this help and exit
```

## Benchmark

DigitalOcean NYC3 eth1 - private network, encryption on

```
# iperf3 -P 1 -R -c 10.247.0.1                                          
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
[  4]   0.00-10.00  sec  1.83 GBytes  1.58 Gbits/sec                  receiver
```
