package components

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type SwapMetric int

const (
	SwapUsed SwapMetric = iota
	SwapUsedPerc
	SwapFree
	SwapFreePerc
)

type Swap struct {
	statusbar.BaseComponentConfig
	metric SwapMetric
}

func NewSwap(metric SwapMetric, interval time.Duration, signal syscall.Signal) *Swap {
	const name = "swap"
	switch metric {
	case SwapUsed:
	case SwapUsedPerc:
	case SwapFree:
	case SwapFreePerc:
	default:
		panic(name + ": unknown metric")
	}
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

	var toFormatted func() string
	switch s.metric {
	case SwapUsed:
		toFormatted = func() string {
			return util.HumanBytes((total - free) * constants.KiB)
		}
	case SwapUsedPerc:
		toFormatted = func() string {
			return fmt.Sprintf("%.1f%%", float64(total-free)/float64(total)*100.0)
		}
	case SwapFree:
		toFormatted = func() string {
			return util.HumanBytes(free * constants.KiB)
		}
	case SwapFreePerc:
		toFormatted = func() string {
			return fmt.Sprintf("%.1f%%", float64(free)/float64(total)*100.0)
		}
	}

	send := func() {
		if err := util.ParseMeminfo(f, buf, fields); err != nil {
			util.Warn("%s: %v", s.Name, err)
			update("")
		} else {
			update(toFormatted())
		}
	}

	send()
	s.BaseComponentConfig.Loop(send, trigger)
}
