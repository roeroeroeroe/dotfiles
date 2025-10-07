package components

import (
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

const diskName = "disk"

func startDisk(cfg statusbar.ComponentConfig, ch chan<- string, trigger <-chan struct{}) {
	name := diskName

	mountpoint, ok := cfg.Arg.(string)
	if !ok || mountpoint == "" {
		util.Warn("%s: Arg not a string or empty, using %s", name, constants.DefaultDiskMount)
		mountpoint = constants.DefaultDiskMount
	}

	var stat unix.Statfs_t
	send := func() {
		err := unix.Statfs(mountpoint, &stat)
		if err != nil {
			util.Warn("%s: statfs %s: %v", name, mountpoint, err)
			ch <- ""
		} else {
			ch <- util.HumanBytes((stat.Blocks - stat.Bfree) * uint64(stat.Frsize))
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
	statusbar.Register(diskName, startDisk)
}
