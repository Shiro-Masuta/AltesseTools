package system

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUMinimalInfo struct {
	Name    string  `json:"name"`    // Nom du CPU
	Cores   int     `json:"cores"`   // Nombre de cores physiques
	Threads int     `json:"threads"` // Nombre de threads logiques
	Mhz     float64 `json:"mhz"`     // Fréquence moyenne
	Percent float64 `json:"percent"` // Utilisation CPU en %
}

func GetCPUMinimalInfo() (*CPUMinimalInfo, error) {
	infos, err := cpu.Info()
	if err != nil || len(infos) == 0 {
		return nil, err
	}

	cpuName := infos[0].ModelName

	// Moyenne des fréquences
	totalMhz := 0.0
	for _, i := range infos {
		totalMhz += i.Mhz
	}
	avgMhz := totalMhz / float64(len(infos))

	// Appels à cpu.Counts() **dans la fonction**
	coresPhysical, _ := cpu.Counts(false) // vrais cœurs
	threadsLogical, _ := cpu.Counts(true) // threads

	percent, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		return nil, err
	}

	return &CPUMinimalInfo{
		Name:    cpuName,
		Cores:   coresPhysical,
		Threads: threadsLogical,
		Mhz:     avgMhz,
		Percent: percent[0],
	}, nil
}
