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
	RedrawDelay = 100 * time.Millisecond
)

// component   description                 argument         note
//
// cpu         perc                        -                -
// disk        used                        mountpoint (/)   -
// ip          -                           iface (eth0)     prefers v4; Interval, Signal unused
// mem         used                        -                -
// net         rx, tx                      iface (eth0)     per max(1, Interval.Seconds())
// sh          -                           command          timeout=Interval*0.75
// swap        used                        -                -
// tcp         ESTABLISHED remote, local   -                -
// text        -                           string           Interval, Signal unused
// time        -                           time.Layout      -
// uptime      -                           -                -

// convenience
const iface = "enp3s0"

// convenience
func text(str string) entry { return entry{"text", cfg{Arg: str}} }

var Components = []entry{
	text("[ "),
	{"sh", cfg{Arg: "player_sb", Interval: time.Second, Signal: 35}},
	text(" vol:"),
	{"sh", cfg{Arg: "volume", Interval: 3 * time.Second, Signal: 36}},
	text(" ]   [ " + iface + " ("),
	{"ip", cfg{Arg: iface}},
	text(") "),
	{"net", cfg{Arg: iface, Interval: time.Second}},
	text(" tcp{"),
	{"tcp", cfg{Interval: time.Second}},
	text("} ]   [ disk:"),
	{"disk", cfg{Arg: "/", Interval: 3 * time.Second}},
	text(" mem:"),
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
