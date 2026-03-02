package components

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"roe/sb/constants"
	"roe/sb/netmon"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type Net struct {
	ifaceName string
	statusbar.BaseComponentConfig
}

func NewNet(ifaceName string, interval time.Duration) *Net {
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

	base := statusbar.NewBaseComponentConfig(name, interval, 0)
	return &Net{ifaceName, *base}
}

func (n *Net) Start(update func(string), _ <-chan struct{}) {
	notify := make(chan struct{}, 1)
	netmon.Get().AddHandler(func(ev netmon.Event) {
		if ev.Type == netmon.LinkEvent && ev.Name == n.ifaceName {
			select {
			case notify <- struct{}{}:
			default:
			}
		}
	})

	buf := make([]byte, constants.NetFileReadBufSize)
	sec := uint64(n.Interval.Seconds())

restart:
	for {
		if _, err := net.InterfaceByName(n.ifaceName); err != nil {
			update("")
			<-notify
			continue
		}

		rxFile, txFile, err := openStatistics(n.ifaceName)
		if err != nil {
			util.Warn("%s: %v", n.Name, err)
			update("")
			<-notify
			continue
		}

		prevRx, err := readU64From(rxFile, buf)
		if err != nil {
			rxFile.Close()
			txFile.Close()
			util.Warn("%s: initial rx: %v", n.Name, err)
			<-notify
			continue
		}
		prevTx, err := readU64From(txFile, buf)
		if err != nil {
			rxFile.Close()
			txFile.Close()
			util.Warn("%s: initial tx: %v", n.Name, err)
			<-notify
			continue
		}

		ticker := time.NewTicker(n.Interval)
		for {
			select {
			case <-ticker.C:
				rx, err := readU64From(rxFile, buf)
				if err != nil {
					util.Warn("%s: %v", n.Name, err)
					ticker.Stop()
					rxFile.Close()
					txFile.Close()
					update("")
					continue restart
				}
				tx, err := readU64From(txFile, buf)
				if err != nil {
					util.Warn("%s: %v", n.Name, err)
					ticker.Stop()
					rxFile.Close()
					txFile.Close()
					update("")
					continue restart
				}

				var rxRate, txRate uint64
				if rx >= prevRx {
					rxRate = (rx - prevRx) / sec
				}
				if tx >= prevTx {
					txRate = (tx - prevTx) / sec
				}

				update(fmt.Sprintf("rx:%s tx:%s",
					util.HumanBytes(rxRate), util.HumanBytes(txRate)))

				prevRx, prevTx = rx, tx
			case <-notify:
				ticker.Stop()
				rxFile.Close()
				txFile.Close()
				update("")
				continue restart
			}
		}
	}
}

func openStatistics(ifaceName string) (*os.File, *os.File, error) {
	statisticsPath := filepath.Join(constants.SysNetClassPath, ifaceName, "statistics")

	rxPath := filepath.Join(statisticsPath, "rx_bytes")
	rxFile, err := os.Open(rxPath)
	if err != nil {
		return nil, nil, err
	}

	txPath := filepath.Join(statisticsPath, "tx_bytes")
	txFile, err := os.Open(txPath)
	if err != nil {
		rxFile.Close()
		return nil, nil, err
	}

	return rxFile, txFile, nil
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
