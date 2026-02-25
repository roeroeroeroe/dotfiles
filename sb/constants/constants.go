package constants

const KiB = 1 << 10

const (
	DefaultCatReadBufSize = 64

	ProcStatPath                = "/proc/stat"
	CPUStatReadBufSize          = 4 * KiB
	ProcStatCPUIdleFieldIndex   = 3
	ProcStatCPUIowaitFieldIndex = 4

	DefaultDiskMount = "/"

	SysBlockPath      = "/sys/block"
	DiskIOReadBufSize = 512

	SysNetClassPath    = "/sys/class/net"
	NetLinkReadBufSize = 8 * KiB

	ProcMeminfoPath    = "/proc/meminfo"
	MemInfoReadBufSize = 2 * KiB

	NetFileReadBufSize = 64

	ProcTCPPath        = "/proc/net/tcp"
	ProcTCP6Path       = "/proc/net/tcp6"
	TCPReadBufSize     = 16 * KiB
	TCPIPDecodeBufSize = 16
	TCPReadChunkSize   = 4 * KiB

	DefaultTimeLayout = "Mon 01/02 15:04:05"

	ProcUptimePath    = "/proc/uptime"
	UptimeReadBufSize = 128
)
