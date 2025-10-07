package components

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"roe/sb/statusbar"
	"roe/sb/util"
)

const shName = "sh"

func startSh(cfg statusbar.ComponentConfig, ch chan<- string, trigger <-chan struct{}) {
	name := shName

	command, ok := cfg.Arg.(string)
	if !ok || command == "" {
		util.Warn("%s: Arg not a string or empty", name)
		ch <- ""
		return
	}

	timeout := cfg.Interval * 3 / 4

	send := func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		out, err := exec.CommandContext(ctx, "sh", "-c", command).Output()
		if err != nil {
			util.Warn("%s: \"%s\" failed: %v", name, command, err)
			ch <- ""
		} else {
			ch <- strings.TrimSpace(string(out))
		}
	}

	send()

	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			send()
		case <-trigger:
			send()
		}
	}
}

func init() {
	statusbar.Register(shName, startSh)
}
