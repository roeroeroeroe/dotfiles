package components

import (
	"net"

	"roe/sb/netmon"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type IPPreference uint8

const (
	IPPrefer4 IPPreference = iota
	IPPrefer6
	IPAny
)

type IP struct {
	ifaceName string
	statusbar.BaseComponentConfig
	ifaceIndex int
	preference IPPreference
}

func NewIP(ifaceName string, preference IPPreference) *IP {
	const name = "ip"
	if ifaceName == "" {
		panic(name + ": empty interface name")
	}

	base := statusbar.NewBaseComponentConfig(name, 0, 0)
	return &IP{
		ifaceName:           ifaceName,
		BaseComponentConfig: *base,
		preference:          preference,
	}
}

func (ip *IP) Start(update func(string), _ <-chan struct{}) {
	send := func() {
		iface, err := net.InterfaceByName(ip.ifaceName)
		if err != nil {
			util.Warn("%s: %v", ip.Name, err)
			ip.ifaceIndex = 0
			update("")
			return
		}
		ip.ifaceIndex = iface.Index

		addrs, err := iface.Addrs()
		if err != nil {
			util.Warn("%s: %v", ip.Name, err)
			update("")
			return
		}
		if len(addrs) == 0 {
			update("")
			return
		}

		var v4, v6 net.IP
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			addr := ipnet.IP
			if addr.IsLoopback() {
				continue
			}

			if addr4 := addr.To4(); addr4 != nil {
				if ip.preference == IPPrefer4 || ip.preference == IPAny {
					update(addr.String())
					return
				}
				if v4 == nil {
					v4 = addr
				}
				continue
			}

			if addr.IsLinkLocalUnicast() {
				continue
			}

			if ip.preference == IPPrefer6 || ip.preference == IPAny {
				update(addr.String())
				return
			}
			if v6 == nil {
				v6 = addr
			}
		}

		var chosen net.IP
		switch ip.preference {
		case IPPrefer4:
			if v4 != nil {
				chosen = v4
			} else {
				chosen = v6
			}
		case IPPrefer6:
			if v6 != nil {
				chosen = v6
			} else {
				chosen = v4
			}
		case IPAny:
			update("")
			return
		}

		if chosen != nil {
			update(chosen.String())
		} else {
			update("")
		}
	}

	send()

	netmon.Get().AddHandler(func(ev netmon.Event) {
		switch ev.Type {
		case netmon.AddrEvent:
			if ip.ifaceIndex != 0 && ev.Index == ip.ifaceIndex {
				send()
			}
		case netmon.LinkEvent:
			if (ip.ifaceIndex != 0 && ev.Index == ip.ifaceIndex) ||
				ev.Name == ip.ifaceName {
				send()
			}
		}
	})

	select {}
}
