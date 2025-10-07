package components

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const memName = "mem"

func parseMeminfo(f *os.File, buf []byte) (uint64, uint64, error) {
	if _, err := f.Seek(0, 0); err != nil {
		return 0, 0, err
	}
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return 0, 0, err
	}
	data := buf[:n]

	var total, available uint64

	for line := range bytes.SplitSeq(data, []byte{'\n'}) {
		if len(line) < 9 {
			continue
		}

		switch {
		case bytes.HasPrefix(line, []byte("MemTotal:")):
			fields := bytes.Fields(line)
			if len(fields) < 2 {
				continue
			}
			num, err := strconv.ParseUint(string(fields[1]), 10, 64)
			if err != nil {
				continue
			}
			total = num * 1024
			if available != 0 {
				return total, available, nil
			}
		case bytes.HasPrefix(line, []byte("MemAvailable:")):
			fields := bytes.Fields(line)
			if len(fields) < 2 {
				continue
			}
			num, err := strconv.ParseUint(string(fields[1]), 10, 64)
			if err != nil {
				continue
			}
			available = num * 1024
			if total != 0 {
				return total, available, nil
			}
		}
	}

	if total == 0 {
		return 0, 0, errors.New("MemTotal not found")
	}

	return total, available, nil
}

func startMem(cfg statusbar.ComponentConfig, ch chan<- string, trigger <-chan struct{}) {
	name := memName

	f, err := os.Open(constants.ProcMeminfoPath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		ch <- ""
		return
	}

	buf := make([]byte, constants.MemInfoReadBufSize)

	send := func() {
		total, available, err := parseMeminfo(f, buf)
		if err != nil {
			util.Warn("%s: %v", name, err)
			ch <- ""
			return
		}
		ch <- util.HumanBytes(total - available)
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
	statusbar.Register(memName, startMem)
}
