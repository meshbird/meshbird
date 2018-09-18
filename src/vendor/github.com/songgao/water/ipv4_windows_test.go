// +build windows
package water

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"testing"
)

func startBroadcast(t *testing.T, dst net.IP) {
	if err := exec.Command("ping", "-n", "2", dst.String()).Start(); err != nil {
		t.Fatal(err)
	}
}

func setupIfce(t *testing.T, ipNet net.IPNet, dev string) {
	sargs := fmt.Sprintf("interface ip set address name=REPLACE_ME source=static addr=REPLACE_ME mask=REPLACE_ME gateway=none")
	args := strings.Split(sargs, " ")
	args[4] = fmt.Sprintf("name=%s", dev)
	args[6] = fmt.Sprintf("addr=%s", ipNet.IP)
	args[7] = fmt.Sprintf("mask=%d.%d.%d.%d", ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3])
	cmd := exec.Command("netsh", args...)
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}
