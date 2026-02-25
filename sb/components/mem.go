package components

import (
	"os"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type Mem struct {
	statusbar.BaseComponentConfig
}

func NewMem(interval time.Duration, signal syscall.Signal) *Mem {
	base := statusbar.NewBaseComponentConfig("mem", interval, signal)
	base.MustBeNonZero()
	return &Mem{*base}
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

	send := func() {
		if err := util.ParseMeminfo(f, buf, fields); err != nil {
			util.Warn("%s: %v", m.Name, err)
			update("")
		} else {
			update(util.HumanBytes((total - available) * constants.KiB))
		}
	}

	send()
	m.BaseComponentConfig.Loop(send, trigger)
}
