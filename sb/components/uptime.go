package components

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const uptimeName = "uptime"

func startUptime(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := uptimeName

	var formatUptime func(string) string
	raw, ok := cfg.Arg.(bool)
	if ok && raw {
		formatUptime = func(uptimeString string) string { return uptimeString }
	} else {
		formatUptime = func(uptimeString string) string {
			fS, err := strconv.ParseFloat(uptimeString, 64)
			if err != nil {
				util.Warn("%s: ParseFloat: %v", name, err)
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
				return fmt.Sprintf("%dm %ds", m, s)
			default:
				return fmt.Sprintf("%ds", s)
			}
		}
	}

	f, err := os.Open(constants.ProcUptimePath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	buf := make([]byte, constants.UptimeReadBufSize)

	send := func() {
		if _, err := f.Seek(0, 0); err != nil {
			util.Warn("%s: seek %s: %v", name, constants.ProcUptimePath, err)
			update("")
			return
		}

		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			util.Warn("%s: read %s: %v", name, constants.ProcUptimePath, err)
			update("")
			return
		}

		index := bytes.IndexByte(buf[:n], ' ')
		if index == -1 {
			util.Warn("%s: unexpected %s format", name, constants.ProcUptimePath)
			update("")
			return
		}

		update(formatUptime(string(buf[:index])))
	}

	send()

	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			send()
		case <-trigger:
			send()
		}
	}
}

func init() {
	statusbar.Register(uptimeName, startUptime)
}
