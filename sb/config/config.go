package config

import (
	"time"

	c "roe/sb/components"
	"roe/sb/statusbar"
)

var (
	// used when a component produces no output
	Placeholder = "n/a"
	// minimum interval between redraws, used to batch multiple rapid updates
	RedrawDelay = 50 * time.Millisecond
)

// convenience
const (
	iface       = "enp3s0"
	blockDevice = "sda"
)

var Components = []statusbar.Component{
	c.NewText("[ "),
	c.NewExec([]string{"player_sb"}, 0, time.Second, 35),
	c.NewText(" vol:"),
	c.NewVolume("muted"),
	c.NewText(" ]   [ "),
	c.NewIP(iface),
	c.NewText(" "),
	c.NewNet(iface, time.Second, 0),
	c.NewText(" tcp{"),
	c.NewTCP(c.TCPModeRemoteAndLocal, time.Second, 0),
	c.NewText("} ]   [ disk{u:"),
	c.NewDisk("/", c.MetricUsed, 3*time.Second, 0),
	c.NewText(" io:"),
	c.NewDiskIO(blockDevice, time.Second, 0),
	c.NewText("} mem:"),
	c.NewMem(c.MetricUsed, 2*time.Second, 0),
	c.NewText(" swap:"),
	c.NewSwap(c.MetricUsed, 2*time.Second, 0),
	c.NewText(" cpu:"),
	c.NewCPU(false, 2*time.Second, 0),
	c.NewText(" ]   [ "),
	c.NewTime("Mon 01/02 15:04:05", time.Second, 0),
	c.NewText(" up:"),
	c.NewUptime(false, 5*time.Second, 0),
	c.NewText(" ]"),
	// a simpler example:
	// c.NewText("cpu:"),
	// c.NewCPU(false, 2*time.Second, 0),
	// c.NewText(", mem:"),
	// c.NewMem(c.MetricUsedPerc, 2*time.Second, 0),
	// c.NewText("   "),
	// c.NewTime("15:04:05", time.Second, 0),
}
