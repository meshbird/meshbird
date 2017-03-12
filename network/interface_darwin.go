// +build darwin

package network

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

/*
#include <unistd.h>
#include <netinet/in.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/kern_control.h>
#include <net/if_utun.h>
#include <sys/ioctl.h>
#include <sys/kern_event.h>

int open_utun(int num) {
	int err;
	int fd;
	struct sockaddr_ctl addr;
	struct ctl_info info;

	fd = socket(PF_SYSTEM, SOCK_DGRAM, SYSPROTO_CONTROL);
	if (fd < 0) {
		return fd;
	}
	memset(&info, 0, sizeof (info));
	strncpy(info.ctl_name, UTUN_CONTROL_NAME, strlen(UTUN_CONTROL_NAME));
	err = ioctl(fd, CTLIOCGINFO, &info);
	if (err < 0) {
		close(fd);
		return err;
	}

	addr.sc_id = info.ctl_id;
	addr.sc_len = sizeof(addr);
	addr.sc_family = AF_SYSTEM;
	addr.ss_sysaddr = AF_SYS_CONTROL;
	addr.sc_unit = num + 1; // utunX where X is sc.sc_unit -1

	err = connect(fd, (struct sockaddr*)&addr, sizeof(addr));
	if (err < 0) {
		// this utun is in use
		close(fd);
		return err;
	}

	return fd;
}
*/
import "C"

type UTUNAccess struct {
	fd int
}

func (a *UTUNAccess) Write(data []byte) (n int, err error) {
	// data from utun is not an ip packet directly. 4 bytes [0, 0, 0, 2] are prepended to it.
 	buf := append([]byte{0, 0, 0, 2}, data...)
	n, err = syscall.Write(a.fd, buf)
	return
}

func (a *UTUNAccess) Read(data []byte) (n int, err error) {
	buf := make([]byte, 1496)
	n, err = syscall.Read(a.fd, buf)
	// data from utun is not an ip packet directly. 4 bytes [0, 0, 0, 2] are prepended to it.
	copy(data, buf[4:])
	return
}

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
	//	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	//	if err != nil {
	//		return nil, err
	//	}
	//	name, err := createInterface(file.Fd(), ifName, cIFF_TAP|cIFF_NO_PI)
	//	if err != nil {
	//		return nil, err
	//	}
	//	ifce = &Interface{isTAP: true, file: file, name: name}
	err = errors.New("unsupported")
	return
}

func newTUN(ifName string) (ifce *Interface, err error) {
	ifce, err = interfaceOpen("tun", "")
	if err != nil {
		return nil, err
	}
	return ifce, nil
}

func createInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
	var req ifReq
	req.Flags = flags
	copy(req.Name[:], ifName)
	createdIFName = strings.Trim(string(req.Name[:]), "\x00")
	return
}

func setPersistent(fd uintptr, persistent bool) error {
	//	var val uintptr = 0
	//	if persistent {
	//		val = 1
	//	}
	//	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETPERSIST), val)
	//	if errno != 0 {
	//		return errno
	//	}
	return nil
}

func interfaceOpen(ifType, ifName string) (*Interface, error) {
	var err error
	if ifType != "tun" && ifType != "tap" {
		return nil, fmt.Errorf("unknown interface type: %s", ifType)
	}
	ifce := new(Interface)
	for i := 0; i < 10; i++ {
		fd := C.open_utun(C.int(i))
		if fd < 0 {
			continue
		}

		ifce.name = fmt.Sprintf("utun%d", i)
		ifce.file = &UTUNAccess{fd: int(fd)}
		break

		//ifPath := fmt.Sprintf("/dev/%s%d", ifType, i)
		//fmt.Println(ifPath)
		//ifce.file, err = os.OpenFile(ifPath, os.O_RDWR, 0)
		//if err != nil {
		//	continue
		//}
		//ifce.name = fmt.Sprintf("%s%d", ifType, i)
		//break
	}
	return ifce, err
}

func AssignIpAddress(iface string, IpAddr string) error {
	ip, ipnet, _ := net.ParseCIDR(IpAddr)
	err := exec.Command("ipconfig", "set", iface, "MANUAL", ip.To4().String(), fmt.Sprintf("0x%s", ipnet.Mask.String())).Run()
	if err != nil {
		return fmt.Errorf("assign ip %s to %s err: %s", IpAddr, iface, err)
	}
	return nil
}

func UpInterface(iface string) error {
	err := exec.Command("ifconfig", iface, "up").Run()
	if err != nil {
		return fmt.Errorf("up interface %s err: %s", iface, err)
	}
	return err
}

func SetMTU(iface string, mtu int) error {
	err := exec.Command("ifconfig", iface, "mtu", strconv.Itoa(mtu)).Run()
	if err != nil {
		return fmt.Errorf("Can't set MTU %s to %s err: %s", iface, strconv.Itoa(mtu), err)
	}
	return nil
}
