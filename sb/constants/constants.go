package constants

const (
	ProcStatPath       = "/proc/stat"
	CPUStatReadBufSize = 4096

	DefaultDiskMount = "/"

	SysBlockPath      = "/sys/block"
	DiskIOReadBufSize = 512

	SysNetClassPath    = "/sys/class/net"
	NetLinkReadBufSize = 8192

	ProcMeminfoPath    = "/proc/meminfo"
	MemInfoReadBufSize = 2048

	NetFileReadBufSize = 64

	ProcTCPPath        = "/proc/net/tcp"
	ProcTCP6Path       = "/proc/net/tcp6"
	TCPReadBufSize     = 16384
	TCPIPDecodeBufSize = 16
	TCPReadChunkSize   = 4096

	DefaultTimeLayout = "Mon 01/02 15:04:05"

	ProcUptimePath    = "/proc/uptime"
	UptimeReadBufSize = 128
)
