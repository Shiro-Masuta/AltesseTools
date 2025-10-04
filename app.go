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

// SelectFiles : ouvre une boîte de dialogue générique avec filtres dynamiques
func (a *App) SelectFiles(title string, filters []runtime.FileFilter, multiple bool) ([]string, error) {
	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   title,
		Filters: filters,
	})
}

func (a *App) SelectFolder(title string) (string, error) {
	if title == "" {
		title = "Sélectionnez un dossier"
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
}

// SaveDroppedFiles : sauvegarde les fichiers droppés dans un dossier temporaire
func (a *App) SaveDroppedFiles(filesData []FileData) ([]string, error) {
	tempDir := filepath.Join(os.TempDir(), "altesse_dropped_files")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("impossible de créer le dossier temporaire: %v", err)
	}

	var savedPaths []string
	for _, fileData := range filesData {
		tempPath := filepath.Join(tempDir, fileData.Name)

		if err := os.WriteFile(tempPath, fileData.Content, 0644); err != nil {
			log.Printf("Erreur sauvegarde %s: %v", fileData.Name, err)
			continue
		}

		savedPaths = append(savedPaths, tempPath)
	}

	return savedPaths, nil
}

// CleanupTempFiles : supprime le dossier temporaire
func (a *App) CleanupTempFiles() error {
	tempDir := filepath.Join(os.TempDir(), "altesse_dropped_files")
	return os.RemoveAll(tempDir)
}
