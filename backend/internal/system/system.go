package system

import "Altesse_Tools_V1.0/backend/pkg"

type SystemMinimalInfo struct {
	CPU    *CPUMinimalInfo    `json:"cpu"`
	Memory *MemoryMinimalInfo `json:"memory"`
	Disk   []*DiskMinimalInfo `json:"disk"`
}

func GetSystemMinimalInfo() (*SystemMinimalInfo, error) {
	results, err := pkg.ParallelRun(
		func() (any, error) { return GetCPUMinimalInfo() },
		func() (any, error) { return GetMemoryMinimalInfo() },
		func() (any, error) { return GetAllDiskInfo() },
	)
	if err != nil {
		return nil, err
	}

	return &SystemMinimalInfo{
		CPU:    results[0].(*CPUMinimalInfo),
		Memory: results[1].(*MemoryMinimalInfo),
		Disk:   results[2].([]*DiskMinimalInfo),
	}, nil
}
