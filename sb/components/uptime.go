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

func startUptime(cfg statusbar.ComponentConfig, ch chan<- string, trigger <-chan struct{}) {
	name := uptimeName

	f, err := os.Open(constants.ProcUptimePath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		ch <- ""
		return
	}

	buf := make([]byte, constants.UptimeReadBufSize)

	send := func() {
		if _, err := f.Seek(0, 0); err != nil {
			util.Warn("%s: seek %s: %v", name, constants.ProcUptimePath, err)
			ch <- ""
			return
		}

		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			util.Warn("%s: read %s: %v", name, constants.ProcUptimePath, err)
			ch <- ""
			return
		}

		fields := bytes.Fields(buf[:n])
		if len(fields) < 1 {
			util.Warn("%s: unexpected %s format", name, constants.ProcUptimePath)
			ch <- ""
			return
		}

		s, err := strconv.ParseFloat(string(fields[0]), 64)
		if err != nil {
			util.Warn("%s: ParseFloat: %v", name, err)
			ch <- ""
			return
		}

		intS := int(s)
		var (
			days    = intS / 86400
			hours   = (intS % 86400) / 3600
			minutes = (intS % 3600) / 60
		)
		switch {
		case days > 0:
			ch <- fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
		case hours > 0:
			ch <- fmt.Sprintf("%dh %dm", hours, minutes)
		case minutes > 0:
			ch <- fmt.Sprintf("%dm", minutes)
		default:
			ch <- fmt.Sprintf("%ds", intS)
		}
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
