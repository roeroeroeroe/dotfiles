package components

import (
	"roe/sb/statusbar"
	"roe/sb/util"

	"golang.org/x/sys/unix"
)

const kernelReleaseName = "kernel_release"

func startKernelRelease(cfg statusbar.ComponentConfig, update func(string), _ <-chan struct{}) {
	name := kernelReleaseName

	utsname := &unix.Utsname{}
	if err := unix.Uname(utsname); err != nil {
		util.Warn("%s: uname: %v", name, err)
		update("")
	} else {
		update(unix.ByteSliceToString(utsname.Release[:]))
	}
}

func init() {
	statusbar.Register(kernelReleaseName, startKernelRelease)
}
