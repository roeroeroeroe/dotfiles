//go:build !with_pulse

package components

import "roe/sb/statusbar"

type Volume struct {
	statusbar.BaseComponentConfig
}

func NewVolume(mutedLabel string) *Volume {
	panic("not implemented")
}

func (v *Volume) Start(_ func(string), _ <-chan struct{}) {
	panic("not implemented")
}
