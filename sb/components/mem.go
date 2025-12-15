package components

import (
	"os"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const memName = "mem"

func startMem(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := memName

	f, err := os.Open(constants.ProcMeminfoPath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	buf := make([]byte, constants.MemInfoReadBufSize)

	var total, available uint64
	fields := []util.MeminfoField{
		{Ptr: &total, Key: []byte("MemTotal:")},
		{Ptr: &available, Key: []byte("MemAvailable:")},
	}

	send := func() {
		if err := util.ParseMeminfo(f, buf, fields); err != nil {
			util.Warn("%s: %v", name, err)
			update("")
		} else {
			update(util.HumanBytes((total - available) * 1024))
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
	statusbar.Register(memName, startMem)
}
