package components

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const cpuName = "cpu"

type cpuSample struct{ idle, total uint64 }

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
			if i == constants.ProcStatCPUIdleFieldIndex ||
				i == constants.ProcStatCPUIowaitFieldIndex {
				idle += num
			}
		}
		return idle, total, nil
	}

	return 0, 0, errors.New("cpu line not found")
}

func parsePerCPUStat(f *os.File, buf []byte, out []cpuSample) ([]cpuSample, error) {
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("empty %s", constants.ProcStatPath)
	}

	b := buf[:n]

	out = out[:0]

	for line := range bytes.SplitSeq(b, []byte{'\n'}) {
		if len(line) < 4 || !bytes.HasPrefix(line, []byte("cpu")) ||
			line[3] < '0' || line[3] > '9' {
			continue
		}

		fields := bytes.Fields(line)[1:]
		if len(fields) < 5 {
			return nil, errors.New("unexpected cpuN line format")
		}

		var idle, total uint64
		for i, field := range fields {
			num, err := util.ParseU64(field)
			if err != nil {
				return nil, fmt.Errorf("parse uint: %v", err)
			}
			total += num
			if i == constants.ProcStatCPUIdleFieldIndex ||
				i == constants.ProcStatCPUIowaitFieldIndex {
				idle += num
			}
		}
		out = append(out, cpuSample{idle, total})
	}

	if len(out) == 0 {
		return nil, errors.New("no cpuN lines found")
	}

	return out, nil
}

func startCPU(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := cpuName

	perCPU := false
	if v, ok := cfg.Arg.(bool); ok && v {
		perCPU = true
	}

	f, err := os.Open(constants.ProcStatPath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	buf := make([]byte, constants.CPUStatReadBufSize)

	var send func()
	if perCPU {
		initCap := runtime.NumCPU()
		var (
			a = make([]cpuSample, 0, initCap)
			b = make([]cpuSample, 0, initCap)
		)

		var (
			prev    = a[:0]
			samples = b[:0]
		)

		outBuf := make([]byte, 0, len("100.0% ")*initCap)

		send = func() {
			var err error
			samples, err = parsePerCPUStat(f, buf, samples[:0])
			if err != nil {
				util.Warn("%s: %v", name, err)
				update("")
				return
			}

			if len(prev) == 0 || len(prev) != len(samples) {
				prev, samples = samples, prev
				return
			}

			outBuf = outBuf[:0]

			for i, s := range samples {
				totalDelta := float64(s.total - prev[i].total)
				if totalDelta == 0 {
					update("")
					prev, samples = samples, prev
					return
				}
				idleDelta := float64(s.idle - prev[i].idle)
				outBuf = strconv.AppendFloat(outBuf, (1.0-idleDelta/totalDelta)*100.0, 'f', 1, 64)
				outBuf = append(outBuf, '%', ' ')
			}
			if len(outBuf) > 0 {
				outBuf = outBuf[:len(outBuf)-1]
			}

			update(string(outBuf))

			prev, samples = samples, prev
		}
	} else {
		var prevIdle, prevTotal uint64
		send = func() {
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
					update(fmt.Sprintf("%.1f%%", (1.0-idleDelta/totalDelta)*100.0))
				}
			}

			prevIdle = idle
			prevTotal = total
		}
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
