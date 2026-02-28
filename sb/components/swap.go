package components

import (
	"syscall"
	"time"

	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
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
	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		util.Warn("%s: sysinfo: %v", s.Name, err)
		update("")
		return
	}
	if info.Totalswap == 0 {
		util.Warn("%s: no swap", s.Name)
		update("")
		return
	}

	send := func() {
		err := unix.Sysinfo(&info)
		if err != nil {
			util.Warn("%s: sysinfo: %v", s.Name, err)
			update("")
		} else {
			var (
				unit = uint64(info.Unit)
				used = (info.Totalswap - info.Freeswap) * unit
			)
			update(s.metric.Format(used, info.Totalswap*unit))
		}
	}

	send()
	s.BaseComponentConfig.Loop(send, trigger)
}
