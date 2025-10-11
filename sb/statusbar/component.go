package statusbar

import (
	"time"

	"roe/sb/util"
)

type ComponentConfig struct {
	Arg      any
	Interval time.Duration
	Signal   int
}

type StartComponent func(cfg ComponentConfig, update func(string), trigger <-chan struct{})

var Registry = make(map[string]StartComponent)

func Register(name string, f StartComponent) {
	if name == "" {
		util.Fatalf("Register: invalid component (empty name)")
	}
	if f == nil {
		util.Fatalf("Register: invalid component (nil StartComponent)")
	}
	if _, exists := Registry[name]; exists {
		util.Fatalf("Register: duplicate name: %s", name)
	}
	Registry[name] = f
}
