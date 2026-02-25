package components

import (
	"os"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type Swap struct {
	statusbar.BaseComponentConfig
}

func NewSwap(interval time.Duration, signal syscall.Signal) *Swap {
	base := statusbar.NewBaseComponentConfig("swap", interval, signal)
	base.MustBeNonZero()
	return &Swap{*base}
}

func (s *Swap) Start(update func(string), trigger <-chan struct{}) {
	f, err := os.Open(constants.ProcMeminfoPath)
	if err != nil {
		util.Warn("%s: %v", s.Name, err)
		update("")
		return
	}

	buf := make([]byte, constants.MemInfoReadBufSize)

	var total, free uint64
	fields := []util.MeminfoField{
		{Ptr: &total, Key: []byte("SwapTotal:")},
		{Ptr: &free, Key: []byte("SwapFree:")},
	}

	send := func() {
		if err := util.ParseMeminfo(f, buf, fields); err != nil {
			util.Warn("%s: %v", s.Name, err)
			update("")
		} else {
			update(util.HumanBytes((total - free) * constants.KiB))
		}
	}

	send()
	s.BaseComponentConfig.Loop(send, trigger)
}
