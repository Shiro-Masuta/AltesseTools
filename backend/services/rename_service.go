package services

import "Altesse_Tools_V1.0/backend/internal/files"

type RenameService struct{}

// Nouvelle instance
func NewRenameService() *RenameService {
	return &RenameService{}
}

// Batch rename
func (r *RenameService) Rename(paths []string, opts files.OptionRename) error {
	return files.RenameBatch(paths, opts)
}
