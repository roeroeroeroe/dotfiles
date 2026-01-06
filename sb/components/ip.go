package components

import (
	"net"
	"syscall"
	"unsafe"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

const ipName = "ip"

func startIP(cfg statusbar.ComponentConfig, update func(string), _ <-chan struct{}) {
	name := ipName

	ifaceName, ok := cfg.Arg.(string)
	if !ok || ifaceName == "" {
		util.Warn("%s: Arg not a string or empty", name)
		update("")
		return
	}
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		util.Warn("%s: interface %s: %v", name, ifaceName, err)
		update("")
		return
	}

	send := func() {
		addrs, err := iface.Addrs()
		if err != nil || len(addrs) == 0 {
			update("")
			return
		}
		var ip string
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
			if ip == "" && ipnet.IP.To16() != nil && !ipnet.IP.IsLinkLocalUnicast() {
				ip = ipnet.IP.String()
			}
		}
		update(ip)
	}

	send()

	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_ROUTE)
	if err != nil {
		util.Warn("%s: socket: %v", name, err)
		return
	}

	if err := unix.Bind(fd, &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: unix.RTMGRP_IPV4_IFADDR | unix.RTMGRP_IPV6_IFADDR,
	}); err != nil {
		unix.Close(fd)
		util.Warn("%s: bind: %v", name, err)
		return
	}

	buf := make([]byte, constants.NetLinkReadBufSize)

	for {
		n, _, err := unix.Recvfrom(fd, buf, 0)
		if err != nil {
			util.Warn("%s: recvfrom: %v", name, err)
			continue
		}

		msgs, err := syscall.ParseNetlinkMessage(buf[:n])
		if err != nil {
			util.Warn("%s: ParseNetlinkMessage: %v", name, err)
			continue
		}

		for _, m := range msgs {
			if (m.Header.Type == unix.RTM_NEWADDR || m.Header.Type == unix.RTM_DELADDR) &&
				len(m.Data) >= unix.SizeofIfAddrmsg &&
				int((*(*unix.IfAddrmsg)(unsafe.Pointer(&m.Data[0]))).Index) == iface.Index {
				send()
				break
			}
		}
	}
}

func init() {
	statusbar.Register(ipName, startIP)
}
