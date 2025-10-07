package components

import (
	"roe/sb/statusbar"
	"roe/sb/util"
)

const textName = "text"

func startText(cfg statusbar.ComponentConfig, ch chan<- string, _ <-chan struct{}) {
	name := textName

	if text, ok := cfg.Arg.(string); !ok || text == "" {
		util.Warn("%s: Arg not a string or empty", name)
		ch <- ""
	} else {
		ch <- text
	}
}

func init() {
	statusbar.Register(textName, startText)
}
