package netmon

import (
	"bytes"
	"sync"
	"syscall"
	"unsafe"

	"roe/sb/constants"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

type EventType uint8

const (
	LinkEvent EventType = iota
	AddrEvent
)

type Event struct {
	Name  string
	Index int
	Type  EventType
}

type Handler func(Event)

type Monitor struct {
	handlers []Handler
	fd       int
	mu       sync.Mutex
	closed   bool
}

var (
	m     *Monitor
	mOnce sync.Once
)

func Get() *Monitor {
	mOnce.Do(func() { m = newMonitor() })
	return m
}

func newMonitor() *Monitor {
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_ROUTE)
	if err != nil {
		util.Warn("netmon: socket: %v", err)
		return &Monitor{closed: true}
	}

	const groups = unix.RTMGRP_LINK | unix.RTMGRP_IPV4_IFADDR | unix.RTMGRP_IPV6_IFADDR
	if err := unix.Bind(fd, &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: groups,
	}); err != nil {
		unix.Close(fd)
		util.Warn("netmon: bind: %v", err)
		return &Monitor{closed: true}
	}

	mon := &Monitor{fd: fd}
	go mon.listen()

	return mon
}

func (m *Monitor) AddHandler(h Handler) {
	if m.closed {
		return
	}

	m.mu.Lock()
	m.handlers = append(m.handlers, h)
	m.mu.Unlock()
}

func (m *Monitor) listen() {
	buf := make([]byte, constants.NetLinkReadBufSize)
	for {
		n, _, err := unix.Recvfrom(m.fd, buf, 0)
		if err != nil {
			if m.closed {
				return
			}
			util.Warn("netmon: recvfrom: %v", err)
			continue
		}
		if n == 0 {
			continue
		}

		msgs, err := syscall.ParseNetlinkMessage(buf[:n])
		if err != nil {
			util.Warn("netmon: ParseNetlinkMessage: %v", err)
			continue
		}

		for _, msg := range msgs {
			m.handleMsg(msg)
		}
	}
}

func (m *Monitor) handleMsg(msg syscall.NetlinkMessage) {
	var ev Event
	switch msg.Header.Type {
	case unix.RTM_NEWLINK, unix.RTM_DELLINK:
		ev.Type = LinkEvent
		if len(msg.Data) < unix.SizeofIfInfomsg {
			return
		}
		ev.Index = int((*unix.IfInfomsg)(unsafe.Pointer(&msg.Data[0])).Index)

		attrs, err := syscall.ParseNetlinkRouteAttr(&msg)
		if err != nil {
			util.Warn("netmon: ParseNetlinkRouteAttr: %v", err)
			break
		}
		for _, a := range attrs {
			if a.Attr.Type == unix.IFLA_IFNAME {
				ev.Name = string(bytes.TrimRight(a.Value, "\x00"))
				break
			}
		}
	case unix.RTM_NEWADDR, unix.RTM_DELADDR:
		ev.Type = AddrEvent
		if len(msg.Data) < unix.SizeofIfAddrmsg {
			return
		}
		ev.Index = int((*unix.IfAddrmsg)(unsafe.Pointer(&msg.Data[0])).Index)
	default:
		return
	}

	m.mu.Lock()
	for _, h := range m.handlers {
		h(ev)
	}
	m.mu.Unlock()
}
