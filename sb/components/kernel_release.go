package components

import (
	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

type KernelRelease struct {
	statusbar.BaseComponentConfig
}

func NewKernelRelease() *KernelRelease {
	base := statusbar.NewBaseComponentConfig("kernel_release", 0, 0)
	return &KernelRelease{*base}
}

func (kr *KernelRelease) Start(update func(string), _ <-chan struct{}) {
	utsname := &unix.Utsname{}
	if err := unix.Uname(utsname); err != nil {
		util.Warn("%s: uname: %v", kr.Name, err)
		update("")
	} else {
		update(unix.ByteSliceToString(utsname.Release[:]))
	}
}
