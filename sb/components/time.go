package components

import (
	"syscall"
	"time"

	"roe/sb/statusbar"
)

type Time struct {
	layout string
	statusbar.BaseComponentConfig
}

func NewTime(layout string, interval time.Duration, signal syscall.Signal) *Time {
	const name = "time"
	if layout == "" {
		panic(name + ": empty layout")
	}

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Time{layout, *base}
}

func (t *Time) Start(update func(string), trigger <-chan struct{}) {
	send := func() { update(time.Now().Format(t.layout)) }
	send()
	t.BaseComponentConfig.Loop(send, trigger)
}
