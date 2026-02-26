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
	metric UsageMetric
}

func NewDisk(mountpoint string, metric UsageMetric, interval time.Duration, signal syscall.Signal) *Disk {
	const name = "disk"
	if mountpoint == "" {
		util.Warn("%s: empty mountpoint, using %s", name, constants.DefaultDiskMount)
		mountpoint = constants.DefaultDiskMount
	}
	metric.Validate(name)

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Disk{mountpoint, *base, metric}
}

func (d *Disk) Start(update func(string), trigger <-chan struct{}) {
	var stat unix.Statfs_t
	send := func() {
		err := unix.Statfs(d.mountpoint, &stat)
		if err != nil {
			util.Warn("%s: statfs %s: %v", d.Name, d.mountpoint, err)
			update("")
		} else {
			// including reserved
			var (
				used  = (stat.Blocks - stat.Bfree) * uint64(stat.Frsize)
				total = stat.Blocks * uint64(stat.Frsize)
			)
			update(d.metric.Format(used, total))
		}
	}

	send()
	d.BaseComponentConfig.Loop(send, trigger)
}
