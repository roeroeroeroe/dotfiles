package components

import (
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

type Disk struct {
	mountpoint string
	statusbar.BaseComponentConfig
}

func NewDisk(mountpoint string, interval time.Duration, signal syscall.Signal) *Disk {
	const name = "disk"
	if mountpoint == "" {
		util.Warn("%s: empty mountpoint, using %s", name, constants.DefaultDiskMount)
		mountpoint = constants.DefaultDiskMount
	}

	base := statusbar.NewBaseComponentConfig(name, interval, signal)
	base.MustBeNonZero()
	return &Disk{mountpoint, *base}
}

func (d *Disk) Start(update func(string), trigger <-chan struct{}) {
	var stat unix.Statfs_t
	send := func() {
		err := unix.Statfs(d.mountpoint, &stat)
		if err != nil {
			util.Warn("%s: statfs %s: %v", d.Name, d.mountpoint, err)
			update("")
		} else {
			update(util.HumanBytes((stat.Blocks - stat.Bfree) * uint64(stat.Frsize)))
		}
	}

	send()
	d.BaseComponentConfig.Loop(send, trigger)
}
