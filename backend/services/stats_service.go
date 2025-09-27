package services

import (
	"fmt"

	"Altesse_Tools_V1.0/backend/internal/stats"
)

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

type WidgetStats struct {
	TotalConverted      int                           `json:"total_converted"`
	Formats             map[string]*stats.FormatStats `json:"formats"`
	TotalSavedInCD      int64                         `json:"total_saved_cd"`
	TotalSavedFloppy    int64                         `json:"total_saved_floppy"`
	TotalOriginalSize   int64                         `json:"total_original_size"`
	TotalCompressedSize int64                         `json:"total_compressed_size"`
}

func (s *StatsService) GetWidgetStats() (*WidgetStats, error) {
	statsData, err := stats.LoadStats()
	if err != nil {
		return nil, fmt.Errorf("impossible de charger les stats: %w", err)
	}

	original, compressed := statsData.TotalSizes() // si tu ajoutes cette méthode

	return &WidgetStats{
		TotalConverted:      statsData.TotalConverted,
		Formats:             statsData.Formats,
		TotalSavedInCD:      statsData.TotalSavedInCD,
		TotalSavedFloppy:    statsData.TotalSavedFloppy,
		TotalOriginalSize:   original,
		TotalCompressedSize: compressed,
	}, nil
}
