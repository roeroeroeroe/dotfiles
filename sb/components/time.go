package components

import (
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const timeName = "time"

func startTime(cfg statusbar.ComponentConfig, ch chan<- string, trigger <-chan struct{}) {
	name := timeName

	layout, ok := cfg.Arg.(string)
	if !ok || layout == "" {
		util.Warn("%s: Arg not a string or empty, using %s", name, constants.DefaultTimeLayout)
		layout = constants.DefaultTimeLayout
	}

	ch <- time.Now().Format(layout)

	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			ch <- time.Now().Format(layout)
		case <-trigger:
			ch <- time.Now().Format(layout)
		}
	}
}

func init() {
	statusbar.Register(timeName, startTime)
}
