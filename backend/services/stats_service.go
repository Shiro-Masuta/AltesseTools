package services

import (
	"fmt"

	"Altesse_Tools_V1.0/backend/pkg"
)

// StatsService permet de récupérer des statistiques pour le frontend
type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

type WidgetStats struct {
	TotalConverted      int                         `json:"total_converted"`
	Formats             map[string]*pkg.FormatStats `json:"formats"`
	TotalSavedInCD      int64                       `json:"total_saved_cd"`
	TotalSavedFloppy    int64                       `json:"total_saved_floppy"`
	TotalOriginalSize   int64                       `json:"total_original_size"`   // en octets
	TotalCompressedSize int64                       `json:"total_compressed_size"` // en octets
}

func (s *StatsService) GetWidgetStats() (*WidgetStats, error) {
	stats, err := pkg.LoadStats()
	if err != nil {
		return nil, fmt.Errorf("impossible de charger les stats: %w", err)
	}

	var totalOriginal int64
	var totalCompressed int64

	for _, f := range stats.Formats {
		totalOriginal += f.OriginalSize
		totalCompressed += f.FinalSize
	}

	return &WidgetStats{
		TotalConverted:      stats.TotalConverted,
		Formats:             stats.Formats,
		TotalSavedInCD:      stats.TotalSavedInCD,
		TotalSavedFloppy:    stats.TotalSavedFloppy,
		TotalOriginalSize:   totalOriginal,
		TotalCompressedSize: totalCompressed,
	}, nil
}
