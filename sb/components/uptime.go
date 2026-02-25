package components

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type Uptime struct {
	statusbar.BaseComponentConfig
	raw bool
}

func NewUptime(raw bool, interval time.Duration, signal syscall.Signal) *Uptime {
	base := statusbar.NewBaseComponentConfig("uptime", interval, signal)
	base.MustBeNonZero()
	return &Uptime{*base, raw}
}

func (u *Uptime) Start(update func(string), trigger <-chan struct{}) {
	var formatUptime func(string) string
	if u.raw {
		formatUptime = func(uptimeString string) string { return uptimeString }
	} else {
		formatUptime = func(uptimeString string) string {
			fS, err := strconv.ParseFloat(uptimeString, 64)
			if err != nil {
				util.Warn("%s: ParseFloat: %v", u.Name, err)
				return ""
			}
			s := int64(fS)
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

	f, err := os.Open(constants.ProcUptimePath)
	if err != nil {
		util.Warn("%s: %v", u.Name, err)
		update("")
		return
	}

	buf := make([]byte, constants.UptimeReadBufSize)

	send := func() {
		if _, err := f.Seek(0, 0); err != nil {
			util.Warn("%s: seek %s: %v", u.Name, constants.ProcUptimePath, err)
			update("")
			return
		}

		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			util.Warn("%s: read %s: %v", u.Name, constants.ProcUptimePath, err)
			update("")
			return
		}

		index := bytes.IndexByte(buf[:n], ' ')
		if index == -1 {
			util.Warn("%s: unexpected %s format", u.Name, constants.ProcUptimePath)
			update("")
			return
		}

		update(formatUptime(string(buf[:index])))
	}

	send()
	u.BaseComponentConfig.Loop(send, trigger)
}
