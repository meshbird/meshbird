// +build linux

package network

import (
	"fmt"
	"os/exec"
	"strconv"
)

const (
	cIFF_TUN   = 0x0001
	cIFF_TAP   = 0x0002
	cIFF_NO_PI = 0x1000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func newTAP(ifName string) (ifce *Interface, err error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	name, err := createInterface(file.Fd(), ifName, cIFF_TAP|cIFF_NO_PI)
	if err != nil {
		return nil, err
	}
	ifce = &Interface{isTAP: true, file: file, name: name}
	return
}

func newTUN(ifName string) (ifce *Interface, err error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	name, err := createInterface(file.Fd(), ifName, cIFF_TUN|cIFF_NO_PI)
	if err != nil {
		return nil, err
	}
	ifce = &Interface{isTAP: false, file: file, name: name}
	return
}

func createInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
	var req ifReq
	req.Flags = flags
	copy(req.Name[:], ifName)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		err = errno
		return
	}
	createdIFName = strings.Trim(string(req.Name[:]), "\x00")
	return
}

func setPersistent(fd uintptr, persistent bool) error {
	var val uintptr = 0
	if persistent {
		val = 1
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETPERSIST), val)
	if errno != 0 {
		return errno
	}
	return nil
}

func interfaceOpen(ifType, ifName string) (*Interface, error) {
	var err error
	if ifType != "tun" && ifType != "tap" {
		return nil, fmt.Errorf("unknown interface type: %s", ifType)
	}
	ifce := new(Interface)
	for i := 0; i < 256; i++ {
		ifPath := fmt.Sprintf("/dev/tun/%s%d", ifType, ifName)
		ifce.file, err = os.OpenFile(ifPath, os.O_RDWR, 0644)
		if err != nil {
			continue
		}
		ifce.name = ifName
	}
	if ifce.file == nil {
		return nil, fmt.Errorf("can't create network interface")
	}
	return ifce, err
}

func AssignIpAddress(iface string, IpAddr string) error {
	err := exec.Command("ip", "addr", "add", IpAddr, "dev", iface).Run()
	if err != nil {
		return fmt.Errorf("assign ip %s to %s err: %s", IpAddr, iface, err)
	}
	return err
}

func UpInterface(iface string) error {
	err := exec.Command("ip", "link", "set", iface, "up").Run()
	if err != nil {
		return fmt.Errorf("up interface %s err: %s", iface, err)
	}
	return err
}

func SetMTU(iface string, mtu int) error {
	err := exec.Command("ip", "link", "set", "mtu", strconv.Itoa(mtu), "dev", iface).Run()
	if err != nil {
		return fmt.Errorf("Can't set MTU %s to %s err: %s", iface, strconv.Itoa(mtu), err)
	}
	return nil
}
