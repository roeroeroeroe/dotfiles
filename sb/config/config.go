package config

import (
	"fmt"
	"time"

	"roe/sb/statusbar"
)

// convenience
type cfg = statusbar.ComponentConfig

type entry struct {
	Name   string
	Config cfg
}

var (
	// used when a component produces no output
	Placeholder = "n/a"
	// minimum interval between redraws, used to batch multiple rapid updates
	RedrawDelay = 50 * time.Millisecond
)

/*
Name             description                 Arg                  Arg type      note

cat              -                           path                 string        reads at most `constants.CatReadBufSize` B
cpu              perc                        perCPU               bool          -
disk             used                        mountpoint (/)       string        -
disk_io          perc                        block device (sda)   string        per max(1, Interval.Milliseconds())
exec             -                           argv                 []string      timeout=Interval*0.75
ip               -                           iface (eth0)         string        prefers v4; Interval, Signal unused
kernel_release   -                           -                    -             Interval, Signal unused
mem              used                        -                    -             -
net              rx, tx                      iface (eth0)         string        per max(1, Interval.Seconds())
swap             used                        -                    -             -
tcp              ESTABLISHED remote, local   -                    -             "local" = connections whose remote address is loopback; "remote" = all other connections
text             -                           text                 string        Interval, Signal unused
time             -                           layout               time.Layout   -
uptime           -                           raw                  bool          -
user             -                           -                    -             Interval, Signal unused
*/

// convenience
const (
	iface       = "enp3s0"
	blockDevice = "sda"
)

// convenience
func text(s string) entry {
	return entry{"text", cfg{Arg: s}}
}

// convenience
func textf(format string, a ...any) entry {
	return entry{"text", cfg{Arg: fmt.Sprintf(format, a...)}}
}

var Components = []entry{
	text("[ "),
	{"exec", cfg{Arg: []string{"player_sb"}, Interval: time.Second, Signal: 35}},
	text(" vol:"),
	{"exec", cfg{Arg: []string{"volume"}, Interval: 3 * time.Second, Signal: 36}},
	text(" ]   [ "),
	{"ip", cfg{Arg: iface}},
	text(" "),
	{"net", cfg{Arg: iface, Interval: time.Second}},
	text(" tcp{"),
	{"tcp", cfg{Interval: time.Second}},
	text("} ]   [ disk{u:"),
	{"disk", cfg{Arg: "/", Interval: 3 * time.Second}},
	text(" io:"),
	{"disk_io", cfg{Arg: blockDevice, Interval: time.Second}},
	text("} mem:"),
	{"mem", cfg{Interval: 2 * time.Second}},
	text(" swap:"),
	{"swap", cfg{Interval: 2 * time.Second}},
	text(" cpu:["),
	{"cpu", cfg{Arg: true, Interval: 2 * time.Second}},
	text("] ]   [ "),
	{"time", cfg{Arg: "Mon 01/02 15:04:05", Interval: time.Second}},
	text(" up:"),
	{"uptime", cfg{Interval: 5 * time.Second}},
	text(" ]"),
	// a simpler example:
	// text("cpu:["),
	// {"cpu", cfg{Arg: true, Interval: 2 * time.Second}},
	// text("], mem:"),
	// {"mem", cfg{Interval: 2 * time.Second}},
	// text("   "),
	// {"time", cfg{Arg: "15:04:05", Interval: time.Second}},
}
