package components

import (
	"fmt"
	"syscall"
	"time"

	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

type Loadavg struct {
	statusbar.BaseComponentConfig
}

func NewLoadavg(interval time.Duration, signal syscall.Signal) *Loadavg {
	base := statusbar.NewBaseComponentConfigNonZero("loadavg", interval, signal)
	return &Loadavg{*base}
}

func (l *Loadavg) Start(update func(string), trigger <-chan struct{}) {
	var info unix.Sysinfo_t
	send := func() {
		if err := unix.Sysinfo(&info); err != nil {
			util.Warn("%s: sysinfo: %v", l.Name, err)
			update("")
		} else {
			const scale = float64(1 << unix.SI_LOAD_SHIFT)
			var (
				load1m  = float64(info.Loads[0]) / scale
				load5m  = float64(info.Loads[1]) / scale
				load15m = float64(info.Loads[2]) / scale
			)
			update(fmt.Sprintf("%.2f, %.2f, %.2f", load1m, load5m, load15m))
		}
	}

	send()
	l.BaseComponentConfig.Loop(send, trigger)
}
