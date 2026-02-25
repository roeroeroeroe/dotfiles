package components

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"roe/sb/statusbar"
	"roe/sb/util"
)

type Exec struct {
	argv []string
	statusbar.BaseComponentConfig
	timeout time.Duration
}

func NewExec(argv []string, timeout, interval time.Duration, signal syscall.Signal) *Exec {
	const name = "exec"
	if len(argv) == 0 {
		panic(name + ": empty argv")
	}
	if timeout <= 0 {
		if timeout < 0 {
			panic(name + ": negative timeout")
		}
		if interval == 0 {
			panic(name + ": interval==0 && timeout==0")
		}
		timeout = interval * 3 / 4
		util.Warn("%s: empty timeout, using %v (interval=%v, %v*0.75=%v)",
			name, timeout, interval, interval, timeout)
	} else if interval != 0 && timeout >= interval {
		panic(name + ": timeout >= interval")
	}

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Exec{argv, *base, timeout}
}

func (e *Exec) Start(update func(string), trigger <-chan struct{}) {
	var (
		eName   = e.argv[0]
		args    = e.argv[1:]
		cmdline = strings.Join(e.argv, " ")
	)

	send := func() {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()

		out, err := exec.CommandContext(ctx, eName, args...).Output()
		if err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				util.Warn("%s: %q timed out after %v", e.Name, cmdline, e.timeout)
			} else {
				util.Warn("%s: %q: %v", e.Name, cmdline, err)
			}
			update("")
		} else {
			update(string(bytes.TrimSpace(out)))
		}
	}

	send()
	e.BaseComponentConfig.Loop(send, trigger)
}
