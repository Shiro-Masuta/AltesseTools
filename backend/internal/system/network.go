package system

import "github.com/shirou/gopsutil/v4/net"

type NetworkMinimalInfo struct {
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
}

func GetNetworkMinimalInfo() (*NetworkMinimalInfo, error) {
	stats, err := net.IOCounters(false) // false = toutes interfaces combinées
	if err != nil || len(stats) == 0 {
		return nil, err
	}

	return &NetworkMinimalInfo{
		BytesSent:   stats[0].BytesSent,
		BytesRecv:   stats[0].BytesRecv,
		PacketsSent: stats[0].PacketsSent,
		PacketsRecv: stats[0].PacketsRecv,
	}, nil
}
