package components

import (
	"bytes"
	"io"
	"os"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type Cat struct {
	path string
	statusbar.BaseComponentConfig
	readBufSize int
	trimSpace   bool
}

func NewCat(path string, readBufSize int, trimSpace bool, interval time.Duration, signal syscall.Signal) *Cat {
	const name = "cat"
	if path == "" {
		panic(name + ": empty path")
	}
	if readBufSize <= 0 {
		if readBufSize < 0 {
			panic(name + ": negative read buffer size")
		}
		util.Warn("%s: read buffer size is 0, using %d", name, constants.DefaultCatReadBufSize)
		readBufSize = constants.DefaultCatReadBufSize
	}

	base := statusbar.NewBaseComponentConfigNonZero(name, interval, signal)
	return &Cat{path, *base, readBufSize, trimSpace}
}

func (c *Cat) Start(update func(string), trigger <-chan struct{}) {
	buf := make([]byte, c.readBufSize)

	var toString func([]byte) string
	if c.trimSpace {
		toString = func(b []byte) string { return string(bytes.TrimSpace(b)) }
	} else {
		toString = func(b []byte) string { return string(b) }
	}

	send := func() {
		f, err := os.Open(c.path)
		if err != nil {
			util.Warn("%s: open %s: %v", c.Name, c.path, err)
			update("")
			return
		}
		defer f.Close()

		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			util.Warn("%s: read %s: %v", c.Name, c.path, err)
			update("")
			return
		}

		if n == len(buf) {
			var tmp [1]byte
			if n, _ := f.Read(tmp[:]); n > 0 {
				util.Warn("%s: truncated output from %s", c.Name, c.path)
			}
		}

		update(toString(buf[:n]))
	}

	send()
	c.BaseComponentConfig.Loop(send, trigger)
}
