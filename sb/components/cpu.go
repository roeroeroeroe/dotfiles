package components

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const cpuName = "cpu"

func parseStat(f *os.File, buf []byte) (uint64, uint64, error) {
	if _, err := f.Seek(0, 0); err != nil {
		return 0, 0, err
	}

	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return 0, 0, err
	}
	if n == 0 {
		return 0, 0, fmt.Errorf("empty %s", constants.ProcStatPath)
	}

	b := buf[:n]

	for line := range bytes.SplitSeq(b, []byte{'\n'}) {
		if len(line) < 4 || !bytes.HasPrefix(line, []byte("cpu ")) {
			continue
		}

		fields := bytes.Fields(line)[1:]
		if len(fields) < 5 {
			return 0, 0, errors.New("unexpected cpu line format")
		}

		var idle, total uint64
		for i, field := range fields {
			num, err := util.ParseU64(field)
			if err != nil {
				return 0, 0, fmt.Errorf("parse uint: %v", err)
			}
			total += num
			if i == 3 || i == 4 {
				idle += num
			}
		}
		return idle, total, nil
	}

	return 0, 0, errors.New("cpu line not found")
}

func startCPU(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := cpuName

	f, err := os.Open(constants.ProcStatPath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	buf := make([]byte, constants.CPUStatReadBufSize)
	var prevIdle, prevTotal uint64

	send := func() {
		idle, total, err := parseStat(f, buf)
		if err != nil {
			util.Warn("%s: %v", name, err)
			update("")
			return
		}

		if prevTotal != 0 {
			totalDelta := float64(total - prevTotal)
			if totalDelta == 0 {
				update("")
			} else {
				idleDelta := float64(idle - prevIdle)
				update(fmt.Sprintf("%.0f%%", (1.0-idleDelta/totalDelta)*100.0))
			}
		}

		prevIdle = idle
		prevTotal = total
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
	statusbar.Register(cpuName, startCPU)
}
