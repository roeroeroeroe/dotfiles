package components

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type DiskIO struct {
	blockDeviceName string
	statusbar.BaseComponentConfig
}

func NewDiskIO(blockDeviceName string, interval time.Duration, signal syscall.Signal) *DiskIO {
	const name = "disk_io"
	if blockDeviceName == "" {
		panic(name + ": empty block device name")
	}
	ms := interval.Milliseconds()
	if ms < 1 {
		panic(name + ": interval < 1ms")
	}

	orig := interval
	interval = time.Duration(ms) * time.Millisecond

	if orig > interval {
		util.Warn("%s: interval adjusted: %v -> %v", name, orig, interval)
	}

	base := statusbar.NewBaseComponentConfig(name, interval, signal)
	return &DiskIO{blockDeviceName, *base}
}

func (d *DiskIO) Start(update func(string), trigger <-chan struct{}) {
	statPath := filepath.Join(constants.SysBlockPath, d.blockDeviceName, "stat")
	file, err := os.Open(statPath)
	if err != nil {
		util.Warn("%s: open %s: %v", d.Name, statPath, err)
		update("")
		return
	}

	var (
		prevWait uint64
		fMs      = float64(d.Interval.Milliseconds())
		buf      = make([]byte, constants.DiskIOReadBufSize)
	)

	send := func() {
		if _, err := file.Seek(0, 0); err != nil {
			util.Warn("%s: seek: %v", d.Name, err)
			update("")
			return
		}

		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			util.Warn("%s: read: %v", d.Name, err)
			update("")
			return
		}
		b := buf[:n]

		if len(b) == 0 {
			util.Warn("%s: empty read", d.Name)
			update("")
			return
		}

		fields := bytes.Fields(b)
		if len(fields) < 10 {
			util.Warn("%s: no io_ticks field", d.Name)
			update("")
			return
		}

		wait, err := util.ParseU64(fields[9])
		if err != nil {
			util.Warn("%s: parse uint: %v", d.Name, err)
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
	d.BaseComponentConfig.Loop(send, trigger)
}
