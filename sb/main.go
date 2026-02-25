package main

import (
	"bytes"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"roe/sb/config"
	"roe/sb/util"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func main() {
	if len(config.Components) == 0 {
		util.Warn("0 components")
		os.Exit(2)
	}

	var (
		bufPool       = sync.Pool{New: func() any { return new(bytes.Buffer) }}
		sigToTriggers = map[syscall.Signal][]chan struct{}{}
		redrawCh      = make(chan struct{}, 1)
	)

	state := make([]atomic.Value, len(config.Components))

	var (
		sFlag bool
		draw  func()
	)
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-s" {
			sFlag = true
			break
		}
	}
	if sFlag {
		draw = func() {
			buf := bufPool.Get().(*bytes.Buffer)
			buf.Reset()

			for i := range len(config.Components) {
				buf.WriteString(state[i].Load().(string))
			}
			buf.WriteByte('\n')

			if _, err := os.Stdout.Write(buf.Bytes()); err != nil {
				util.Warn("Stdout.Write: %v", err)
			}
			bufPool.Put(buf)
		}
	} else {
		conn, err := xgb.NewConn()
		if err != nil {
			util.Warn("x conn: %v", err)
			os.Exit(1)
		}
		root := xproto.Setup(conn).DefaultScreen(conn).Root

		draw = func() {
			buf := bufPool.Get().(*bytes.Buffer)
			buf.Reset()

			for i := range len(config.Components) {
				buf.WriteString(state[i].Load().(string))
			}

			b := buf.Bytes()
			if err := xproto.ChangePropertyChecked(
				conn,
				xproto.PropModeReplace,
				root,
				xproto.AtomWmName,
				xproto.AtomString,
				8,
				uint32(len(b)),
				b,
			).Check(); err != nil {
				util.Warn("set name: %v", err)
			}
			bufPool.Put(buf)
		}
	}

	noopTrigger := make(chan struct{}, 1)

	for i, c := range config.Components {
		state[i].Store(config.Placeholder)

		var trigger chan struct{}
		if s := c.Signal(); s != 0 {
			trigger = make(chan struct{}, 1)
			sig := syscall.Signal(s)
			sigToTriggers[sig] = append(sigToTriggers[sig], trigger)
		} else {
			trigger = noopTrigger
		}

		update := func(str string) {
			if str == "" {
				str = config.Placeholder
			}
			if state[i].Load().(string) != str {
				state[i].Store(str)
				select {
				case redrawCh <- struct{}{}:
				default:
				}
			}
		}

		go c.Start(update, trigger)
	}

	redrawTimer := time.NewTimer(config.RedrawDelay)
	redrawTimer.Stop()
	go func() {
		for {
			select {
			case <-redrawCh:
				if !redrawTimer.Stop() {
					select {
					case <-redrawTimer.C:
					default:
					}
				}
				redrawTimer.Reset(config.RedrawDelay)
			case <-redrawTimer.C:
				draw()
			}
		}
	}()

	if len(sigToTriggers) > 0 {
		sigCh := make(chan os.Signal, 1)
		var sigs []os.Signal
		for s := range sigToTriggers {
			sigs = append(sigs, s)
		}
		signal.Notify(sigCh, sigs...)

		for sig := range sigCh {
			for _, trig := range sigToTriggers[sig.(syscall.Signal)] {
				select {
				case trig <- struct{}{}:
				default:
					util.Warn("signal %d already pending", sig)
				}
			}
		}
	} else {
		select {}
	}
}
