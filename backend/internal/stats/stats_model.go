package stats

type FormatStats struct {
	Count        int   `json:"count"`         // nombre d'images converties
	OriginalSize int64 `json:"original_size"` // taille totale avant conversion
	FinalSize    int64 `json:"final_size"`    // taille totale après conversion
}

type Stats struct {
	TotalConverted   int                     `json:"total_converted"`
	Formats          map[string]*FormatStats `json:"formats"`
	TotalSavedInCD   int64                   `json:"total_saved_cd"`     // en CD 700 Mo
	TotalSavedFloppy int64                   `json:"total_saved_floppy"` // en disquettes 1.44 Mo
}
