package util

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"roe/sb/constants"
)

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
