package components

import (
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const timeName = "time"

func startTime(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := timeName

	layout, ok := cfg.Arg.(string)
	if !ok || layout == "" {
		util.Warn("%s: Arg not a string or empty, using %s", name, constants.DefaultTimeLayout)
		layout = constants.DefaultTimeLayout
	}

	update(time.Now().Format(layout))

	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			update(time.Now().Format(layout))
		case <-trigger:
			update(time.Now().Format(layout))
		}
	}
}

func init() {
	statusbar.Register(timeName, startTime)
}
