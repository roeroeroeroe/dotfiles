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

		fields := bytes.Fields(buf[:n])
		if len(fields) < 1 {
			util.Warn("%s: unexpected %s format", name, constants.ProcUptimePath)
			update("")
			return
		}

		s, err := strconv.ParseFloat(string(fields[0]), 64)
		if err != nil {
			util.Warn("%s: ParseFloat: %v", name, err)
			update("")
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
			update(fmt.Sprintf("%dd %dh %dm", days, hours, minutes))
		case hours > 0:
			update(fmt.Sprintf("%dh %dm", hours, minutes))
		case minutes > 0:
			update(fmt.Sprintf("%dm", minutes))
		default:
			update(fmt.Sprintf("%ds", intS))
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
