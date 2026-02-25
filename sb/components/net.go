package components

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type Net struct {
	ifaceName string
	statusbar.BaseComponentConfig
}

func NewNet(ifaceName string, interval time.Duration, signal syscall.Signal) *Net {
	const name = "net"
	if ifaceName == "" {
		panic(name + ": empty interface name")
	}
	sec := uint64(interval.Seconds())
	if sec < 1 {
		panic(name + ": interval < 1s")
	}

	orig := interval
	interval = time.Duration(sec) * time.Second

	if orig > interval {
		util.Warn("%s: interval adjusted: %v -> %v", name, orig, interval)
	}

	base := statusbar.NewBaseComponentConfig(name, interval, signal)
	return &Net{ifaceName, *base}
}

func (n *Net) Start(update func(string), trigger <-chan struct{}) {
	iface, err := net.InterfaceByName(n.ifaceName)
	if err != nil {
		util.Warn("%s: interface %s: %v", n.Name, n.ifaceName, err)
		update("")
		return
	}

	rxPath := filepath.Join(constants.SysNetClassPath, iface.Name, "statistics", "rx_bytes")
	rxFile, err := os.Open(rxPath)
	if err != nil {
		util.Warn("%s: open %s: %v", n.Name, rxPath, err)
		update("")
		return
	}

	txPath := filepath.Join(constants.SysNetClassPath, iface.Name, "statistics", "tx_bytes")
	txFile, err := os.Open(txPath)
	if err != nil {
		rxFile.Close()
		util.Warn("%s: open %s: %v", n.Name, txPath, err)
		update("")
		return
	}

	buf := make([]byte, constants.NetFileReadBufSize)

	prevRx, err := readU64From(rxFile, buf)
	if err != nil {
		rxFile.Close()
		txFile.Close()
		util.Warn("%s: read %s: %v", n.Name, rxPath, err)
		update("")
		return
	}
	prevTx, err := readU64From(txFile, buf)
	if err != nil {
		rxFile.Close()
		txFile.Close()
		util.Warn("%s: read %s: %v", n.Name, txPath, err)
		update("")
		return
	}

	sec := uint64(n.Interval.Seconds())

	send := func() {
		rx, err := readU64From(rxFile, buf)
		if err != nil {
			util.Warn("%s: read %s: %v", n.Name, rxPath, err)
			update("")
			return
		}
		tx, err := readU64From(txFile, buf)
		if err != nil {
			util.Warn("%s: read %s: %v", n.Name, txPath, err)
			update("")
			return
		}

		update(fmt.Sprintf("rx:%s tx:%s",
			util.HumanBytes((rx-prevRx)/sec), util.HumanBytes((tx-prevTx)/sec)))

		prevRx = rx
		prevTx = tx
	}

	send()
	n.BaseComponentConfig.Loop(send, trigger)
}

func readU64From(f *os.File, buf []byte) (uint64, error) {
	if _, err := f.Seek(0, 0); err != nil {
		return 0, err
	}

	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	if n == 0 {
		return 0, errors.New("empty")
	}

	return util.ParseU64(buf[:n])
}
