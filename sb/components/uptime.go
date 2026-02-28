package components

import (
	"fmt"
	"strconv"
	"syscall"
	"time"

	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

type Uptime struct {
	statusbar.BaseComponentConfig
	raw bool
}

func NewUptime(raw bool, interval time.Duration, signal syscall.Signal) *Uptime {
	base := statusbar.NewBaseComponentConfigNonZero("uptime", interval, signal)
	return &Uptime{*base, raw}
}

func (u *Uptime) Start(update func(string), trigger <-chan struct{}) {
	var formatUptime func(float64) string
	if u.raw {
		formatUptime = func(upSec float64) string {
			return strconv.FormatFloat(upSec, 'f', 0, 64)
		}
	} else {
		formatUptime = func(upSec float64) string {
			s := int64(upSec)
			const (
				minute = 60
				hour   = 60 * minute
				day    = 24 * hour
				week   = 7 * day
			)
			var (
				w = s / week
				d = (s % week) / day
				h = (s % day) / hour
				m = (s % hour) / minute
			)
			switch {
			case w > 0:
				return fmt.Sprintf("%dw %dd %dh", w, d, h)
			case d > 0:
				return fmt.Sprintf("%dd %dh %dm", d, h, m)
			case h > 0:
				return fmt.Sprintf("%dh %dm", h, m)
			case m > 0:
				return fmt.Sprintf("%dm %ds", m, s%minute)
			default:
				return fmt.Sprintf("%ds", s)
			}
		}
	}

	var info unix.Sysinfo_t
	send := func() {
		if err := unix.Sysinfo(&info); err != nil {
			util.Warn("%s: sysinfo: %v", u.Name, err)
			update("")
		} else {
			update(formatUptime(float64(info.Uptime)))
		}
	}

	send()
	u.BaseComponentConfig.Loop(send, trigger)
}
