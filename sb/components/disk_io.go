package components

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const diskIOName = "disk_io"

func startDiskIO(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := diskIOName

	blockDeviceName, ok := cfg.Arg.(string)
	if !ok || blockDeviceName == "" {
		util.Warn("%s: Arg not a string or empty", name)
		update("")
		return
	}

	statPath := filepath.Join(constants.SysBlockPath, blockDeviceName, "stat")
	file, err := os.Open(statPath)
	if err != nil {
		util.Warn("%s: open %s: %v", name, statPath, err)
		update("")
		return
	}

	var (
		prevWait uint64
		ms       = max(1, uint64(cfg.Interval.Milliseconds()))
		fMs      = float64(ms)
		buf      = make([]byte, constants.DiskIOReadBufSize)
	)

	send := func() {
		if _, err := file.Seek(0, 0); err != nil {
			util.Warn("%s: seek: %v", name, err)
			update("")
			return
		}

		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			util.Warn("%s: read: %v", name, err)
			update("")
			return
		}
		b := buf[:n]

		if len(b) == 0 {
			util.Warn("%s: empty read", name)
			update("")
			return
		}

		fields := bytes.Fields(b)
		if len(fields) < 10 {
			util.Warn("%s: no io_ticks field", name)
			update("")
			return
		}

		wait, err := util.ParseU64(fields[9])
		if err != nil {
			util.Warn("%s: parse uint: %v", name, err)
			update("")
			return
		}
		if prevWait == 0 {
			prevWait = wait
			update("")
			return
		}

		update(fmt.Sprintf("%.1f%%", (float64(wait-prevWait)/fMs)*100.0))
		prevWait = wait
	}

	send()

	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
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
	statusbar.Register(diskIOName, startDiskIO)
}
