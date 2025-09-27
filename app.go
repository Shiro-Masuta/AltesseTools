package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type FileData struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

func (a *App) SelectImagePaths() ([]string, error) {
	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Sélectionnez des images",
		Filters: []runtime.FileFilter{
			{DisplayName: "Images", Pattern: "*.png;*.jpg;*.jpeg;*.webp;*.avif;*.bmp;*.tiff"},
		},
	})
}

// SelectOutputFolder ouvre un sélecteur de dossier pour choisir le dossier de sortie
func (a *App) SelectOutputFolder() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Sélectionnez un dossier de sortie",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (a *App) HandleDroppedFiles(paths []string) ([]string, error) {
	valid := []string{}
	for _, p := range paths {
		if fileExists(p) && isImage(p) {
			valid = append(valid, p)
		}
	}
	return valid, nil
}

// Nouvelle fonction pour sauvegarder les fichiers droppés
func (a *App) SaveDroppedFiles(filesData []FileData) ([]string, error) {
	// Créer un dossier temporaire
	tempDir := filepath.Join(os.TempDir(), "altesse_dropped_files")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("impossible de créer le dossier temporaire: %v", err)
	}

	var validPaths []string

	for _, fileData := range filesData {
		// Créer le chemin temporaire
		tempPath := filepath.Join(tempDir, fileData.Name)

		// Sauvegarder le fichier
		if err := os.WriteFile(tempPath, fileData.Content, 0644); err != nil {
			log.Printf("Erreur sauvegarde %s: %v", fileData.Name, err)
			continue
		}

		// Vérifier si c'est une image valide
		if isImage(tempPath) {
			validPaths = append(validPaths, tempPath)
		} else {
			// Supprimer le fichier si ce n'est pas une image
			os.Remove(tempPath)
		}
	}

	return validPaths, nil
}

// Fonction pour nettoyer les fichiers temporaires (optionnel)
func (a *App) CleanupTempFiles() error {
	tempDir := filepath.Join(os.TempDir(), "altesse_dropped_files")
	return os.RemoveAll(tempDir)
}

// Vérifie si le fichier existe
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Vérifie si l’extension correspond à une image
func isImage(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".avif", ".bmp", ".tiff":
		return true
	}
	return false
}
