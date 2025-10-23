package components

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"roe/sb/statusbar"
	"roe/sb/util"
)

const execName = "exec"

func startExec(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := execName

	argv, ok := cfg.Arg.([]string)
	if !ok || len(argv) == 0 || argv[0] == "" {
		util.Warn("%s: Arg not a []string or empty", name)
		update("")
		return
	}

	var (
		eName   = argv[0]
		args    = argv[1:]
		cmdline = strings.Join(argv, " ")
	)

	timeout := cfg.Interval * 3 / 4

	send := func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		out, err := exec.CommandContext(ctx, eName, args...).Output()
		if err != nil {
			util.Warn("%s: \"%s\" failed: %v", name, cmdline, err)
			update("")
		} else {
			update(string(bytes.TrimSpace(out)))
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
	statusbar.Register(execName, startExec)
}
