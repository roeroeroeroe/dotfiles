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

func startNet(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := netName

	iface, err := util.ArgOrFirstUpIface(cfg.Arg)
	if err != nil {
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	var (
		rxPath = filepath.Join(constants.SysNetClassPath, iface.Name, "statistics", "rx_bytes")
		txPath = filepath.Join(constants.SysNetClassPath, iface.Name, "statistics", "tx_bytes")
	)

	rxFile, err := os.Open(rxPath)
	if err != nil {
		util.Warn("%s: open %s: %v", name, rxPath, err)
		update("")
		return
	}

	txFile, err := os.Open(txPath)
	if err != nil {
		rxFile.Close()
		util.Warn("%s: open %s: %v", name, txPath, err)
		update("")
		return
	}

	buf := make([]byte, constants.NetFileReadBufSize)

	prevRx, err := readU64From(rxFile, buf)
	if err != nil {
		util.Warn("%s: read %s: %v", name, rxPath, err)
		update("")
		return
	}
	prevTx, err := readU64From(txFile, buf)
	if err != nil {
		util.Warn("%s: read %s: %v", name, txPath, err)
		update("")
		return
	}

	sec := max(1, uint64(cfg.Interval.Seconds()))

	send := func() {
		rx, err := readU64From(rxFile, buf)
		if err != nil {
			util.Warn("%s: read rx: %v", name, err)
			update("")
			return
		}
		tx, err := readU64From(txFile, buf)
		if err != nil {
			util.Warn("%s: read tx: %v", name, err)
			update("")
			return
		}

		update(fmt.Sprintf("rx:%s tx:%s",
			util.HumanBytes((rx-prevRx)/sec), util.HumanBytes((tx-prevTx)/sec)))

		prevRx = rx
		prevTx = tx
	}

	send()

	ticker := time.NewTicker(time.Duration(sec) * time.Second)
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
