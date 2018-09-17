// +build linux darwin

package water

import (
	"net"
	"os/exec"
	"testing"
)

func startBroadcast(t *testing.T, dst net.IP) {
	if err := exec.Command("ping", "-b", "-c", "2", dst.String()).Start(); err != nil {
		t.Fatal(err)
	}
}

func setupIfce(t *testing.T, ipNet net.IPNet, dev string) {
	if err := exec.Command("ip", "link", "set", dev, "up").Run(); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("ip", "addr", "add", ipNet.String(), "dev", dev).Run(); err != nil {
		t.Fatal(err)
	}
}
