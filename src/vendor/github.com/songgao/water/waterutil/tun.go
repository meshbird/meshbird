package waterutil

func IsIPv4(packet []byte) bool {
	return 4 == (packet[0] >> 4)
}

func IsIPv6(packet []byte) bool {
	return 6 == (packet[0] >> 4)
}
