package statusbar

import (
	"syscall"
	"time"
)

type Component interface {
	Start(update func(string), trigger <-chan struct{})
	Signal() syscall.Signal
}

type BaseComponentConfig struct {
	Name     string
	Interval time.Duration
	Sig      syscall.Signal
}

func NewBaseComponentConfig(name string, interval time.Duration, signal syscall.Signal) *BaseComponentConfig {
	if interval < 0 {
		panic(name + ": negative interval")
	}
	return &BaseComponentConfig{name, interval, signal}
}

func (b BaseComponentConfig) MustBeNonZero() {
	if b.Interval == 0 && b.Sig == 0 {
		panic(b.Name + ": both interval and signal are 0")
	}
}

func (b BaseComponentConfig) Signal() syscall.Signal {
	return b.Sig
}

func (b BaseComponentConfig) Loop(send func(), trigger <-chan struct{}) {
	if b.Interval == 0 {
		for {
			<-trigger
			send()
		}
	}

	ticker := time.NewTicker(b.Interval)

	if b.Sig == 0 {
		for {
			<-ticker.C
			send()
		}
	}

	for {
		select {
		case <-ticker.C:
			send()
		case <-trigger:
			send()
		}
	}
}
