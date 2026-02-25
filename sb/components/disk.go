package components

import (
	"fmt"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

type DiskMetric int

const (
	DiskUsed DiskMetric = iota
	DiskUsedPerc
	DiskFree
	DiskFreePerc
)

type Disk struct {
	mountpoint string
	statusbar.BaseComponentConfig
	metric DiskMetric
}

func NewDisk(mountpoint string, metric DiskMetric, interval time.Duration, signal syscall.Signal) *Disk {
	const name = "disk"
	if mountpoint == "" {
		util.Warn("%s: empty mountpoint, using %s", name, constants.DefaultDiskMount)
		mountpoint = constants.DefaultDiskMount
	}
	switch metric {
	case DiskUsed:
	case DiskUsedPerc:
	case DiskFree:
	case DiskFreePerc:
	default:
		panic(name + ": unknown metric")
	}

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Disk{mountpoint, *base, metric}
}

func (d *Disk) Start(update func(string), trigger <-chan struct{}) {
	var (
		stat        unix.Statfs_t
		toFormatted func() string
	)
	switch d.metric {
	case DiskUsed:
		toFormatted = func() string {
			return util.HumanBytes((stat.Blocks - stat.Bfree) * uint64(stat.Frsize))
		}
	case DiskUsedPerc:
		toFormatted = func() string {
			return fmt.Sprintf("%.1f%%", float64(stat.Blocks-stat.Bfree)/float64(stat.Blocks)*100.0)
		}
	case DiskFree:
		toFormatted = func() string {
			return util.HumanBytes(stat.Bavail * uint64(stat.Frsize))
		}
	case DiskFreePerc:
		toFormatted = func() string {
			return fmt.Sprintf("%.1f%%", float64(stat.Bavail)/float64(stat.Blocks)*100.0)
		}
	}
	send := func() {
		err := unix.Statfs(d.mountpoint, &stat)
		if err != nil {
			util.Warn("%s: statfs %s: %v", d.Name, d.mountpoint, err)
			update("")
		} else {
			update(toFormatted())
		}
	}

	send()
	d.BaseComponentConfig.Loop(send, trigger)
}
