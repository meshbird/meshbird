# Meshbird 

Meshbird - distributed private networking. [Twitter](https://twitter.com/meshbird), [Website](http://meshbird.com/)


![MeshBird](https://avatars0.githubusercontent.com/u/16837838?v=3&u=dbd30ffcc7383854dba5a66425ce9fe0591b03ac&s=700)

## Intro

Meshbird create distributed private networking between servers, containers, virtual machines and any computers in different datacenters, different countries, different cloud providers. All traffic transmit directly to recepient peer without passing any gateways. Meshbird do not require any centralized servers. Meshbird is absolutly decentralized distributed private networking.

For example, user can create private network between DigitalOcean's droplets in each datacenter and link it together by executing one command. All traffic will be encrypted with strong AES-256.

## Supported OS

- Linux (386, amd64, arm)
- Darwin (386, amd64)

## Demo of ssh connection via our secure tunnel

Demo install and setup Meshbird on DigitalOcean droplets and home laptop for creating distributed private network and provide access to ssh server over Meshbird network.

![SSHDemo](https://raw.githubusercontent.com/meshbird/meshbird/master/demos/ssh_demo.gif)

Full video can be found at https://www.youtube.com/watch?v=sW5ZIcfX7w8

[Other demos](http://meshbird.com/post/demos/)

## Technologies used

1. DHT - this is our strongest side that open way to build fully distributed and secured private networking
2. STUN
3. AES-256 - traffic encription
4. uTP node communication

## Roadmap

1. Better encryption
2. Windows/OSx support
3. IPv6
4. Much more
5. IP Load Balancing
6. HTTP Load Balancing
7. DNS Service Discovery
8. Additional peer discovery plugins

## Quick Start

### Install

```bash
$ curl http://meshbird.com/install.sh | sh
```

or if you have Go compiler 

```bash
$ go get github.com/meshbird/meshbird
```

### Generate new network secret key

```bash
$ meshbird new
```

### Join to the private network with interface tunX

```bash
$ MESHBIRD_KEY="<key>" meshbird join
```

## Ideas

Any ideas feel free to send to email: miolini@gmail.com.
