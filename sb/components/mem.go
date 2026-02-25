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

type MemMetric int

const (
	MemUsed MemMetric = iota
	MemUsedPerc
	MemFree
	MemFreePerc
)

type Mem struct {
	statusbar.BaseComponentConfig
	metric MemMetric
}

func NewMem(metric MemMetric, interval time.Duration, signal syscall.Signal) *Mem {
	const name = "mem"
	switch metric {
	case MemUsed:
	case MemUsedPerc:
	case MemFree:
	case MemFreePerc:
	default:
		panic(name + ": unknown metric")
	}

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Mem{*base, metric}
}

func (m *Mem) Start(update func(string), trigger <-chan struct{}) {
	f, err := os.Open(constants.ProcMeminfoPath)
	if err != nil {
		util.Warn("%s: %v", m.Name, err)
		update("")
		return
	}

	buf := make([]byte, constants.MemInfoReadBufSize)

	var total, available uint64
	fields := []util.MeminfoField{
		{Ptr: &total, Key: []byte("MemTotal:")},
		{Ptr: &available, Key: []byte("MemAvailable:")},
	}

	var toFormatted func() string
	switch m.metric {
	case MemUsed:
		toFormatted = func() string {
			return util.HumanBytes((total - available) * constants.KiB)
		}
	case MemUsedPerc:
		toFormatted = func() string {
			return fmt.Sprintf("%.1f%%", float64(total-available)/float64(total)*100.0)
		}
	case MemFree:
		toFormatted = func() string {
			return util.HumanBytes(available * constants.KiB)
		}
	case MemFreePerc:
		toFormatted = func() string {
			return fmt.Sprintf("%.1f%%", float64(available)/float64(total)*100.0)
		}
	}

	send := func() {
		if err := util.ParseMeminfo(f, buf, fields); err != nil {
			util.Warn("%s: %v", m.Name, err)
			update("")
		} else {
			update(toFormatted())
		}
	}

	send()
	m.BaseComponentConfig.Loop(send, trigger)
}
