package config

import (
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

// component        description                 argument             note
//
// cat              -                           path                 reads at most `constants.CatReadBufSize` B
// cpu              perc                        -                    -
// disk             used                        mountpoint (/)       -
// disk_io          perc                        block device (sda)   per max(1, Interval.Milliseconds())
// exec             -                           argv []string        timeout=Interval*0.75
// ip               -                           iface (eth0)         prefers v4; Interval, Signal unused
// kernel_release   -                           -                    Interval, Signal unused
// mem              used                        -                    -
// net              rx, tx                      iface (eth0)         per max(1, Interval.Seconds())
// swap             used                        -                    -
// tcp              ESTABLISHED remote, local   -                    -
// text             -                           string               Interval, Signal unused
// time             -                           time.Layout          -
// uptime           -                           -                    -
// user             -                           -                    Interval, Signal unused

// convenience
const iface = "enp3s0"
const blockDevice = "sda"

// convenience
func text(str string) entry { return entry{"text", cfg{Arg: str}} }

var Components = []entry{
	text("[ "),
	{"exec", cfg{Arg: []string{"player_sb"}, Interval: time.Second, Signal: 35}},
	text(" vol:"),
	{"exec", cfg{Arg: []string{"volume"}, Interval: 3 * time.Second, Signal: 36}},
	text(" ]   [ " + iface + " ("),
	{"ip", cfg{Arg: iface}},
	text(") "),
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
	text(" cpu:"),
	{"cpu", cfg{Interval: 2 * time.Second}},
	text(" ]   [ "),
	{"time", cfg{Arg: "Mon 01/02 15:04:05", Interval: time.Second}},
	text(" up:"),
	{"uptime", cfg{Interval: 5 * time.Second}},
	text(" ]"),
	// a simpler example:
	// text("cpu:"),
	// {"cpu", cfg{Interval: 2 * time.Second}},
	// text(", mem:"),
	// {"mem", cfg{Interval: 2 * time.Second}},
	// text("   "),
	// {"time", cfg{Arg: "15:04:05", Interval: time.Second}},
}
