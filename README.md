# Meshbird 

Meshbird - distributed private networking. [Twitter](https://twitter.com/meshbird), [Website](http://meshbird.com/)


![MeshBird](https://avatars0.githubusercontent.com/u/16837838?v=3&u=dbd30ffcc7383854dba5a66425ce9fe0591b03ac&s=700)

## Intro
Meshbird enables distributed private networking across geographically dispersed datacenters. With Meshbird senders can send data directly to recipients without the need for gateways or centralized servers. Meshbird is compatible across countries, cloud providers, and container technologies.

**Meshbird encrypts all traffic, but currently uses AES-256 in unauthenticated CBC mode, which does not guarantee confidentiality. Meshbird's encryption is currently not suited for working with sensitive data, unless an additional layer of protection like SSH is used.**

For example, with a single command, Meshbird can be used to create a private network connecting a Docker container in DigitalOcean's in New York data center to a laptop in a Kyrgyzstan café.

## Supported OS

- Linux (386, amd64, arm)
- Darwin (386, amd64)

## Demo of ssh connection via our secure tunnel

Demo install and setup Meshbird on DigitalOcean droplets and home laptop for creating distributed private network and provide access to ssh server over Meshbird network.

![SSHDemo](https://raw.githubusercontent.com/meshbird/meshbird/master/demos/ssh_demo.gif)

Full video can be found at https://www.youtube.com/watch?v=sW5ZIcfX7w8

[Other demos](http://meshbird.com/post/demos/)

## Technologies used

1. **DHT** - distributed node discovery (based on seed IP address list)
2. **STUN** - NAT traversal
3. **AES-256** - traffic encryption
4. **uTP** - node communication with low-latency congestion control

## Roadmap

1. Better encryption
2. Windows support
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
