package util

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

type MeminfoField struct {
	Ptr   *uint64
	Key   []byte
	found bool
}

func Warn(format string, a ...any)   { fmt.Fprintf(os.Stderr, format+"\n", a...) }
func Fatalf(format string, a ...any) { Warn(format, a...); os.Exit(1) }

var iecPrefixes = [...]string{
	"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB",
}

func HumanBytes(b uint64) string {
	const (
		base  = 1024
		fBase = float64(base)
	)
	if b < base {
		return fmt.Sprintf("%d %s", b, iecPrefixes[0])
	}

	var (
		n = float64(b)
		i = 0
	)
	for ; i+1 < len(iecPrefixes) && n >= fBase; i++ {
		n /= base
	}

	return fmt.Sprintf("%.1f %s", n, iecPrefixes[i])
}

func ParseU64(buf []byte) (uint64, error) {
	i := 0
	for i < len(buf) &&
		(buf[i] == ' ' || buf[i] == '\t' || buf[i] == '\n' || buf[i] == '\r') {
		i++
	}
	if i == len(buf) {
		return 0, errors.New("no digits")
	}
	var (
		num   uint64
		found bool
	)
	for ; i < len(buf); i++ {
		c := buf[i]
		if c < '0' || c > '9' {
			break
		}
		found = true
		num = num*10 + uint64(c-'0')
	}
	if !found {
		return 0, errors.New("no digits")
	}
	return num, nil
}

func ParseMeminfo(f *os.File, buf []byte, fields []MeminfoField) error {
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return err
	}
	b := buf[:n]

	found := 0
	for i := range fields {
		fields[i].found = false
	}

	start := 0
	for start < len(b) {
		end := bytes.IndexByte(b[start:], '\n')
		if end < 0 {
			end = len(b) - start
		}
		line := b[start : start+end]
		start += end + 1

		for i := range fields {
			field := &fields[i]
			if field.found || !bytes.HasPrefix(line, field.Key) {
				continue
			}

			num, err := ParseU64(line[len(field.Key):])
			if err != nil {
				return err
			}

			*field.Ptr = num
			found++
			if found == len(fields) {
				return nil
			}
			field.found = true
			break
		}
	}

	return fmt.Errorf("missing %d fields", len(fields)-found)
}
