package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	statsFile string
	mutex     sync.Mutex
)

// init initialise le chemin de statsFile dans Documents/AltesseTools
func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		statsFile = "stats.json" // fallback si pas de home
		return
	}

	docs := filepath.Join(home, "Documents", "AltesseTools")
	_ = os.MkdirAll(docs, 0755) // créer le dossier si inexistant
	statsFile = filepath.Join(docs, "stats.json")
}

func LoadStats() (*Stats, error) {
	file, err := os.ReadFile(statsFile)
	if err != nil {
		return &Stats{
			Formats: make(map[string]*FormatStats),
		}, nil
	}

	var s Stats
	if err := json.Unmarshal(file, &s); err != nil {
		return nil, err
	}

	if s.Formats == nil {
		s.Formats = make(map[string]*FormatStats)
	}
	return &s, nil
}

func SaveStats(s *Stats) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statsFile, data, 0644)
}

// RegisterConversion enregistre une conversion et met à jour les statistiques
func RegisterConversion(format string, originalSize, finalSize int64) error {
	mutex.Lock()
	defer mutex.Unlock()

	stats, err := LoadStats()
	if err != nil {
		return err
	}

	stats.TotalConverted++

	if _, exists := stats.Formats[format]; !exists {
		stats.Formats[format] = &FormatStats{}
	}

	stats.Formats[format].Count++
	stats.Formats[format].OriginalSize += originalSize
	stats.Formats[format].FinalSize += finalSize

	// Calcul des économies
	totalSaved := int64(0)
	for _, f := range stats.Formats {
		totalSaved += f.OriginalSize - f.FinalSize
	}
	stats.TotalSavedInCD = totalSaved / (700 * 1024 * 1024) // CD 700 Mo
	stats.TotalSavedFloppy = totalSaved / 1474560

	return SaveStats(stats)
}

func (s *Stats) TotalSizes() (original, compressed int64) {
	for _, f := range s.Formats {
		original += f.OriginalSize
		compressed += f.FinalSize
	}
	return
}
