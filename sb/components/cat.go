package components

import (
	"io"
	"os"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const catName = "cat"

func startCat(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := catName

	path, ok := cfg.Arg.(string)
	if !ok || path == "" {
		util.Warn("%s: Arg not a string or empty", name)
		update("")
		return
	}

	buf := make([]byte, constants.CatReadBufSize)

	send := func() {
		f, err := os.Open(path)
		if err != nil {
			util.Warn("%s: open %s: %v", name, path, err)
			update("")
			return
		}
		defer f.Close()

		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			util.Warn("%s: read %s: %v", name, path, err)
			update("")
			return
		}

		if n == len(buf) {
			var tmp [1]byte
			if n, _ := f.Read(tmp[:]); n > 0 {
				util.Warn("%s: truncated output from %s", name, path)
			}
		}

		update(string(buf[:n]))
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
	statusbar.Register(catName, startCat)
}
