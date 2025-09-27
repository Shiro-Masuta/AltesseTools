package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sync"

	"Altesse_Tools_V1.0/backend/internal/images"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ConverterService struct {
	ctx context.Context
}

func NewConverterService() *ConverterService {
	return &ConverterService{}
}

func (c *ConverterService) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *ConverterService) ConvertManyFromFilesToBase64Parallel(paths []string, opts *images.Options) ([]string, error) {
	encoded := make([]string, len(paths))

	var mu sync.Mutex
	converted := 0

	_, err := images.ParallelConvert(paths, func(i int, path string) ([]byte, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("échec ouverture '%s' : %w", path, err)
		}
		defer file.Close()

		// Conversion de l'image en mémoire
		data, err := images.ConvertFromReader(file, opts)
		if err != nil {
			return nil, err
		}

		// Stockage en base64 dans le tableau
		encoded[i] = base64.StdEncoding.EncodeToString(data)

		// Incrémenter le compteur de fichiers traités
		mu.Lock()
		converted++
		current := converted
		mu.Unlock()

		// Émission de la progression au frontend via Wails
		runtime.EventsEmit(c.ctx, "conversion-progress", map[string]interface{}{
			"current": current,
			"total":   len(paths),
			"path":    path,
		})

		return data, nil
	})

	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func (c *ConverterService) ConvertManyToFolder(paths []string, outputDir string, opts *images.Options) ([]string, error) {
	// Vérifie que le dossier existe (sinon le crée)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("impossible de créer le dossier de sortie : %w", err)
	}

	outputPaths := make([]string, len(paths))
	format := opts.Format

	// Compteur thread-safe pour progression
	var mu sync.Mutex
	converted := 0

	// Conversion en parallèle
	_, err := images.ParallelConvert(paths, func(i int, path string) ([]byte, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("échec ouverture '%s' : %w", path, err)
		}
		defer file.Close()

		// Conversion de l'image
		data, err := images.ConvertFromReader(file, opts)
		if err != nil {
			return nil, err
		}

		// Écriture immédiate du fichier converti
		base := filepath.Base(path)
		name := base[:len(base)-len(filepath.Ext(base))]
		outPath := filepath.Join(outputDir, fmt.Sprintf("%s_converted.%s", name, format))
		if err := os.WriteFile(outPath, data, 0o644); err != nil {
			return nil, fmt.Errorf("échec écriture '%s' : %w", outPath, err)
		}
		outputPaths[i] = outPath

		// Incrémenter le compteur de fichiers convertis
		mu.Lock()
		converted++
		current := converted
		mu.Unlock()

		// Émission de la progression au frontend via Wails
		runtime.EventsEmit(c.ctx, "conversion-progress", map[string]interface{}{
			"current": current,
			"total":   len(paths),
			"path":    path,
			"output":  outPath,
		})

		return data, nil
	})

	if err != nil {
		return nil, err
	}

	return outputPaths, nil
}

func (c *ConverterService) GetImagePreview(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Déterminer le type MIME
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))

	// Retourner en base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
