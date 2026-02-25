package components

import (
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type Time struct {
	layout string
	statusbar.BaseComponentConfig
}

func NewTime(layout string, interval time.Duration, signal syscall.Signal) *Time {
	const name = "time"
	if layout == "" {
		util.Warn("%s: empty layout, using %s", name, constants.DefaultTimeLayout)
		layout = constants.DefaultTimeLayout
	}

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Time{layout, *base}
}

func (t *Time) Start(update func(string), trigger <-chan struct{}) {
	send := func() { update(time.Now().Format(t.layout)) }
	send()
	t.BaseComponentConfig.Loop(send, trigger)
}
