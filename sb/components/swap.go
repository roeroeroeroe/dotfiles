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
	metric UsageMetric
}

func NewSwap(metric UsageMetric, interval time.Duration, signal syscall.Signal) *Swap {
	const name = "swap"
	metric.Validate(name)

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Swap{*base, metric}
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

	if err := util.ParseMeminfo(f, buf, fields); err != nil {
		util.Warn("%s: initial read: %v", s.Name, err)
		update("")
		return
	}
	if total == 0 {
		util.Warn("%s: no swap", s.Name)
		update("")
		return
	}

	send := func() {
		if err := util.ParseMeminfo(f, buf, fields); err != nil {
			util.Warn("%s: %v", s.Name, err)
			update("")
		} else {
			used := (total - free) * constants.KiB
			update(s.metric.Format(used, total*constants.KiB))
		}
	}

	send()
	s.BaseComponentConfig.Loop(send, trigger)
}
