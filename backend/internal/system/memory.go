package system

import "github.com/shirou/gopsutil/v4/mem"

type MemoryMinimalInfo struct {
	Total   uint64  `json:"total"`   // RAM totale en bytes
	Used    uint64  `json:"used"`    // RAM utilisée
	Free    uint64  `json:"free"`    // RAM libre
	Percent float64 `json:"percent"` // % utilisé
}

func GetMemoryMinimalInfo() (*MemoryMinimalInfo, error) {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryMinimalInfo{
		Total:   vMem.Total, // RAM totale installée
		Used:    vMem.Used,  // RAM utilisée
		Free:    vMem.Free,  // RAM libre
		Percent: vMem.UsedPercent,
	}, nil
}
