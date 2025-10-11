package main

import (
	"bytes"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "roe/sb/components"
	"roe/sb/config"
	"roe/sb/statusbar"
	"roe/sb/util"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func main() {
	var (
		triggers      []chan struct{}
		mu            sync.RWMutex
		state         = make([]string, 0, len(config.Components))
		bufPool       = sync.Pool{New: func() any { return new(bytes.Buffer) }}
		sigToTriggers = map[syscall.Signal][]chan struct{}{}
		redrawCh      = make(chan struct{}, 1)
	)

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

			mu.RLock()
			for i := range state {
				buf.WriteString(state[i])
			}
			mu.RUnlock()

			buf.WriteByte('\n')
			if _, err := os.Stdout.Write(buf.Bytes()); err != nil {
				util.Warn("Stdout.Write: %v", err)
			}

			bufPool.Put(buf)
		}
	} else {
		conn, err := xgb.NewConn()
		if err != nil {
			util.Fatalf("x conn: %v", err)
		}
		root := xproto.Setup(conn).DefaultScreen(conn).Root

		draw = func() {
			buf := bufPool.Get().(*bytes.Buffer)
			buf.Reset()

			mu.RLock()
			for i := range state {
				buf.WriteString(state[i])
			}
			mu.RUnlock()

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

	for _, e := range config.Components {
		component, ok := statusbar.Registry[e.Name]
		if !ok {
			util.Warn("unknown component %s, skipping", e.Name)
			continue
		}

		index := len(state)
		state = append(state, config.Placeholder)

		var trigger chan struct{}
		if s := e.Config.Signal; s != 0 {
			trigger = make(chan struct{}, 1)
			sig := syscall.Signal(s)
			sigToTriggers[sig] = append(sigToTriggers[sig], trigger)
		} else {
			trigger = noopTrigger
		}
		triggers = append(triggers, trigger)

		update := func(str string) {
			if str == "" {
				str = config.Placeholder
			}
			mu.Lock()
			state[index] = str
			mu.Unlock()

			select {
			case redrawCh <- struct{}{}:
			default:
			}
		}

		go component(e.Config, update, trigger)
	}

	if len(state) == 0 {
		util.Fatalf("no active components")
	}

	redrawTimer := time.AfterFunc(config.RedrawDelay, draw)
	go func() {
		for range redrawCh {
			redrawTimer.Stop()
			redrawTimer = time.AfterFunc(config.RedrawDelay, draw)
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
