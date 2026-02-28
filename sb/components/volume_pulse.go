//go:build with_pulse

//go:generate go run ../cmd/gen_volume_pulse_constants

package components

/*
#cgo pkg-config: libpulse
#cgo CFLAGS: -Wall -Werror -Wextra
#include <stdint.h>
#include <stdlib.h>

typedef struct volume_pulse volume_pulse_t;

volume_pulse_t *volume_pulse_new(const char *, uintptr_t);
void volume_pulse_run_mainloop(volume_pulse_t *);
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"roe/sb/statusbar"
)

type Volume struct {
	statusbar.BaseComponentConfig
	mutedLabel string
	update     func(string)
	vol        *C.volume_pulse_t
	h          cgo.Handle
}

//export update
func update(str *C.char, h C.uintptr_t) {
	v := cgo.Handle(h).Value().(*Volume)
	v.update(C.GoString(str))
}

func NewVolume(mutedLabel string) *Volume {
	const name = "volume"
	if mutedLabel == "" {
		panic(name + ": empty muted label")
	}

	base := statusbar.NewBaseComponentConfig(name, 0, 0)
	return &Volume{BaseComponentConfig: *base, mutedLabel: mutedLabel}
}

func (v *Volume) Start(update func(string), _ <-chan struct{}) {
	v.update = update
	v.h = cgo.NewHandle(v)

	mutedLabel := C.CString(v.mutedLabel)
	v.vol = C.volume_pulse_new(mutedLabel, C.uintptr_t(v.h))
	if v.vol == nil {
		panic(v.Name + ": volume_pulse_new failed")
	}
	C.free(unsafe.Pointer(mutedLabel))

	go C.volume_pulse_run_mainloop(v.vol)
}
