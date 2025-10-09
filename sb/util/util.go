package util

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"roe/sb/constants"
)

type MeminfoField struct {
	Key   []byte
	Ptr   *uint64
	found bool
}

func Warn(format string, args ...any)   { fmt.Fprintf(os.Stderr, format+"\n", args...) }
func Fatalf(format string, args ...any) { Warn(format, args...); os.Exit(1) }

var iecPrefixes = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

const (
	base  = 1024
	fBase = float64(base)
)

func HumanBytes(b uint64) string {
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

func ArgOrFirstUpIface(arg any) (*net.Interface, error) {
	if name, ok := arg.(string); ok && name != "" {
		iface, err := net.InterfaceByName(name)
		if err != nil {
			return nil, fmt.Errorf("interface %s: %v", name, err)
		}
		return iface, nil
	}
	ents, err := os.ReadDir(constants.SysNetClassPath)
	if err != nil {
		return nil, err
	}
	for _, e := range ents {
		name := e.Name()
		if name == "lo" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(constants.SysNetClassPath, name, "operstate"))
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(b)) == "up" {
			return net.InterfaceByName(name)
		}
	}
	return nil, errors.New("no up interface")
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
