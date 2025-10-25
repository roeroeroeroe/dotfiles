package main

import (
	"bytes"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
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
		bufPool       = sync.Pool{New: func() any { return new(bytes.Buffer) }}
		sigToTriggers = map[syscall.Signal][]chan struct{}{}
		redrawCh      = make(chan struct{}, 1)
	)

	activeCount := 0
	for _, e := range config.Components {
		if _, ok := statusbar.Registry[e.Name]; ok {
			activeCount++
		} else {
			util.Warn("unknown component %s, skipping", e.Name)
		}
	}
	if activeCount == 0 {
		util.Fatalf("no active components")
	}
	state := make([]atomic.Value, activeCount)

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

			for i := range activeCount {
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
			util.Fatalf("x conn: %v", err)
		}
		root := xproto.Setup(conn).DefaultScreen(conn).Root

		draw = func() {
			buf := bufPool.Get().(*bytes.Buffer)
			buf.Reset()

			for i := range activeCount {
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

	slotIndex := 0
	for _, e := range config.Components {
		component, ok := statusbar.Registry[e.Name]
		if !ok {
			continue
		}

		index := slotIndex
		state[index].Store(config.Placeholder)

		var trigger chan struct{}
		if s := e.Config.Signal; s != 0 {
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
			if state[index].Load().(string) != str {
				state[index].Store(str)
				select {
				case redrawCh <- struct{}{}:
				default:
				}
			}
		}

		go component(e.Config, update, trigger)
		slotIndex++
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
