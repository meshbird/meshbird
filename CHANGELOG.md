# Changelog

## v2.0

- trasport changed to TCP
- encryption changed to AES-256-GCM
- removed DHT
- removed NAT traversal
- transport may create many TCP connections to avoid big RTT related problems
- each node may connect via local private network in region if target node located 
  in same region
- bootstraping and peer exchange from one or more nodes

## v0.2
- add OS X support
- improve network performance
- UPnP forked, translated to English and error checking improvement
- minor changes