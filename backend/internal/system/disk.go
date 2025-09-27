package system

import "github.com/shirou/gopsutil/v4/disk"

type DiskMinimalInfo struct {
	Mountpoint string  `json:"mountpoint"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Percent    float64 `json:"percent"`
}

func GetAllDiskInfo() ([]*DiskMinimalInfo, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	var disks []*DiskMinimalInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // ignorer partitions non accessibles
		}
		disks = append(disks, &DiskMinimalInfo{
			Mountpoint: p.Mountpoint,
			Total:      usage.Total,
			Used:       usage.Used,
			Percent:    usage.UsedPercent,
		})
	}
	return disks, nil
}
