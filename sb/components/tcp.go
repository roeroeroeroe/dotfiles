package components

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

type TCP struct {
	statusbar.BaseComponentConfig
}

func NewTCP(interval time.Duration, signal syscall.Signal) *TCP {
	base := statusbar.NewBaseComponentConfigNonZero("tcp", interval, signal)
	return &TCP{*base}
}

func (tcp *TCP) Start(update func(string), trigger <-chan struct{}) {
	var (
		files = make([]*os.File, 2)
		err   error
	)
	files[0], err = os.Open(constants.ProcTCPPath)
	if err != nil {
		util.Warn("%s: %v", tcp.Name, err)
		update("")
		return
	}

	files[1], err = os.Open(constants.ProcTCP6Path)
	if err != nil {
		files[0].Close()
		util.Warn("%s: %v", tcp.Name, err)
		update("")
		return
	}

	var (
		bufs = [][]byte{
			make([]byte, 0, constants.TCPReadBufSize),
			make([]byte, 0, constants.TCPReadBufSize),
		}
		chunks = [][]byte{
			make([]byte, constants.TCPReadChunkSize),
			make([]byte, constants.TCPReadChunkSize),
		}
		decodeBuf = make([]byte, constants.TCPIPDecodeBufSize)
	)

	send := func() {
		remote, local, err := parseTCP(files, bufs, chunks, decodeBuf)
		if err != nil {
			util.Warn("%s: %v", tcp.Name, err)
			update("")
		} else {
			update(fmt.Sprintf("r:%d l:%d", remote, local))
		}
	}

	send()
	tcp.BaseComponentConfig.Loop(send, trigger)
}

func parseTCP(files []*os.File, bufs, chunks [][]byte, decodeBuf []byte) (uint, uint, error) {
	var remote, local uint

	for i, f := range files {
		if _, err := f.Seek(0, 0); err != nil {
			return 0, 0, err
		}

		var (
			buf   = bufs[i][:0]
			chunk = chunks[i]
		)
		for {
			n, err := f.Read(chunk)
			if n > 0 {
				buf = append(buf, chunk[:n]...)
			}
			if err != nil {
				if err != io.EOF {
					return 0, 0, err
				}
				break
			}
		}
		bufs[i] = buf

		if len(buf) == 0 {
			continue
		}

		j := bytes.IndexByte(buf, '\n')
		if j < 0 || j+1 >= len(buf) {
			continue
		}
		j++

		var need int
		if i == 0 {
			need = net.IPv4len
		} else {
			need = net.IPv6len
		}

		for j < len(buf) {
			lf := bytes.IndexByte(buf[j:], '\n')
			var line []byte
			if lf < 0 {
				line = buf[j:]
				j = len(buf)
			} else {
				line = buf[j : j+lf]
				j += lf + 1
			}
			if len(line) == 0 {
				continue
			}

			fields := bytes.Fields(line)
			if len(fields) < 4 ||
				len(fields[3]) != 2 || fields[3][0] != '0' || fields[3][1] != '1' {
				continue
			}

			addrField := fields[2]
			colon := bytes.IndexByte(addrField, ':')
			if colon <= 0 {
				continue
			}
			ipHex := addrField[:colon]

			if len(ipHex) != need*2 {
				continue
			}

			dst := decodeBuf[:need]
			if _, err := hex.Decode(dst, ipHex); err != nil {
				continue
			}

			if need == net.IPv4len {
				dst[0], dst[1], dst[2], dst[3] =
					dst[3], dst[2], dst[1], dst[0]
			} else {
				for k := 0; k < 16; k += 4 {
					dst[k], dst[k+1], dst[k+2], dst[k+3] =
						dst[k+3], dst[k+2], dst[k+1], dst[k]
				}
			}

			if net.IP(dst).IsLoopback() {
				local++
			} else {
				remote++
			}
		}
	}

	return remote, local, nil
}
