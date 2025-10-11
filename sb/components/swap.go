package components

import (
	"os"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const swapName = "swap"

func startSwap(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := swapName

	f, err := os.Open(constants.ProcMeminfoPath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	buf := make([]byte, constants.MemInfoReadBufSize)

	var total, free uint64
	fields := []util.MeminfoField{
		{Key: []byte("SwapTotal:"), Ptr: &total},
		{Key: []byte("SwapFree:"), Ptr: &free},
	}

	send := func() {
		if err := util.ParseMeminfo(f, buf, fields); err != nil {
			util.Warn("%s: %v", name, err)
			update("")
		} else {
			update(util.HumanBytes((total - free) * 1024))
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
	statusbar.Register(swapName, startSwap)
}
