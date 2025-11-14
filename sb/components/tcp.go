package components

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"roe/sb/constants"
	"roe/sb/statusbar"
	"roe/sb/util"
)

const tcpName = "tcp"

func parseTCP(files []*os.File, bufs, chunks [][]byte, decodeBuf []byte) (uint, uint) {
	var remote, local uint

	for i, f := range files {
		if _, err := f.Seek(0, 0); err != nil {
			util.Warn("%s: seek: %v", tcpName, err)
			continue
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
					util.Warn("%s: read: %v", tcpName, err)
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

			var need int
			if i == 1 {
				need = 16
			} else {
				need = 4
			}
			if len(ipHex) != need*2 {
				continue
			}

			dst := decodeBuf[:need]
			if _, err := hex.Decode(dst, ipHex); err != nil {
				continue
			}
			if need == 16 {
				for k := 0; k < 16; k += 4 {
					dst[k], dst[k+1], dst[k+2], dst[k+3] =
						dst[k+3], dst[k+2], dst[k+1], dst[k]
				}
			} else {
				dst[0], dst[1], dst[2], dst[3] =
					dst[3], dst[2], dst[1], dst[0]
			}

			if net.IP(dst).IsLoopback() {
				local++
			} else {
				remote++
			}
		}
	}

	return remote, local
}

func startTCP(cfg statusbar.ComponentConfig, update func(string), trigger <-chan struct{}) {
	name := tcpName

	var (
		files = make([]*os.File, 2)
		err   error
	)
	files[0], err = os.Open(constants.ProcTCPPath)
	if err != nil {
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	files[1], err = os.Open(constants.ProcTCP6Path)
	if err != nil {
		files[0].Close()
		util.Warn("%s: %v", name, err)
		update("")
		return
	}

	var (
		bufs      = [][]byte{make([]byte, 0, constants.TCPReadBufSize), make([]byte, 0, constants.TCPReadBufSize)}
		chunks    = [][]byte{make([]byte, constants.TCPReadChunkSize), make([]byte, constants.TCPReadChunkSize)}
		decodeBuf = make([]byte, constants.TCPIPDecodeBufSize)
	)

	send := func() {
		remote, local := parseTCP(files, bufs, chunks, decodeBuf)
		update(fmt.Sprintf("r:%d l:%d", remote, local))
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
	statusbar.Register(tcpName, startTCP)
}
