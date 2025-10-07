package components

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const netName = "net"

func readUintFrom(f *os.File, buf []byte) (uint64, error) {
	n, err := f.ReadAt(buf, 0)
	if err != nil && err != io.EOF {
		return 0, err
	}
	if n == 0 {
		return 0, errors.New("empty")
	}
	b := buf[:n]

	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n' || b[i] == '\r') {
		i++
	}
	if i == len(b) {
		return 0, errors.New("no digits")
	}

	var (
		found bool
		num   uint64
	)
	for _, c := range b[i:] {
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

func startNet(cfg statusbar.ComponentConfig, ch chan<- string, trigger <-chan struct{}) {
	name := netName

	iface, err := util.ArgOrFirstUpIface(cfg.Arg)
	if err != nil {
		util.Warn("%s: %v", name, err)
		ch <- ""
		return
	}

	var (
		rxPath = filepath.Join(constants.SysNetClassPath, iface.Name, "statistics", "rx_bytes")
		txPath = filepath.Join(constants.SysNetClassPath, iface.Name, "statistics", "tx_bytes")
	)

	rxFile, err := os.Open(rxPath)
	if err != nil {
		util.Warn("%s: open %s: %v", name, rxPath, err)
		ch <- ""
		return
	}

	txFile, err := os.Open(txPath)
	if err != nil {
		rxFile.Close()
		util.Warn("%s: open %s: %v", name, txPath, err)
		ch <- ""
		return
	}

	buf := make([]byte, constants.NetFileReadBufSize)

	prevRx, err := readUintFrom(rxFile, buf)
	if err != nil {
		util.Warn("%s: read %s: %v", name, rxPath, err)
		ch <- ""
		return
	}
	prevTx, err := readUintFrom(txFile, buf)
	if err != nil {
		util.Warn("%s: read %s: %v", name, txPath, err)
		ch <- ""
		return
	}

	sec := max(1, uint64(cfg.Interval.Seconds()))

	send := func() {
		rx, err := readUintFrom(rxFile, buf)
		if err != nil {
			util.Warn("%s: read rx: %v", name, err)
			ch <- ""
			return
		}
		tx, err := readUintFrom(txFile, buf)
		if err != nil {
			util.Warn("%s: read tx: %v", name, err)
			ch <- ""
			return
		}

		ch <- fmt.Sprintf("rx:%s tx:%s",
			util.HumanBytes((rx-prevRx)/sec), util.HumanBytes((tx-prevTx)/sec))

		prevRx = rx
		prevTx = tx
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
	statusbar.Register(netName, startNet)
}
