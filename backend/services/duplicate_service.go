package services

import "Altesse_Tools_V1.0/backend/internal/files"

// DuplicateService fournit des méthodes pour rechercher et supprimer des doublons
type DuplicateService struct{}

func NewDuplicateService() *DuplicateService {
	return &DuplicateService{}
}

// Search recherche les doublons dans un dossier
func (s *DuplicateService) Search(root string) (map[string][]string, error) {
	return files.FindDuplicates(root)
}

// Delete supprime les doublons trouvés (en parallèle)
func (s *DuplicateService) Delete(duplicates map[string][]string) ([]string, error) {
	return files.DeleteDuplicates(duplicates)
}
