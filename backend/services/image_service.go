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
	"Altesse_Tools_V1.0/backend/internal/stats"
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

func (c *ConverterService) recordStats(originalSize, finalSize int64, format string) {
	err := stats.RegisterConversion(format, originalSize, finalSize)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("Erreur enregistrement stats: %v", err))
	}
}

func (c *ConverterService) ConvertManyFromFilesToBase64Parallel(paths []string, opts *images.Options) ([]string, error) {
	encoded := make([]string, len(paths))

	var mu sync.Mutex
	converted := 0

	_, err := images.ParallelConvert(paths, func(i int, path string) ([]byte, error) {
		// Taille avant conversion
		info, _ := os.Stat(path)
		originalSize := info.Size()

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("échec ouverture '%s' : %w", path, err)
		}
		defer file.Close()

		// Conversion
		data, err := images.ConvertFromReader(file, opts)
		if err != nil {
			return nil, err
		}

		// Taille finale
		finalSize := int64(len(data))
		c.recordStats(originalSize, finalSize, opts.Format)

		// Stockage en base64
		encoded[i] = base64.StdEncoding.EncodeToString(data)

		mu.Lock()
		converted++
		current := converted
		mu.Unlock()

		// Event vers le frontend
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
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("impossible de créer le dossier de sortie : %w", err)
	}

	outputPaths := make([]string, len(paths))
	format := opts.Format

	var mu sync.Mutex
	converted := 0

	_, err := images.ParallelConvert(paths, func(i int, path string) ([]byte, error) {
		// Taille avant conversion
		info, _ := os.Stat(path)
		originalSize := info.Size()

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("échec ouverture '%s' : %w", path, err)
		}
		defer file.Close()

		// Conversion
		data, err := images.ConvertFromReader(file, opts)
		if err != nil {
			return nil, err
		}

		// Taille finale
		finalSize := int64(len(data))
		c.recordStats(originalSize, finalSize, format)

		// Écriture du fichier
		base := filepath.Base(path)
		name := base[:len(base)-len(filepath.Ext(base))]
		outPath := filepath.Join(outputDir, fmt.Sprintf("%s_converted.%s", name, format))
		if err := os.WriteFile(outPath, data, 0o644); err != nil {
			return nil, fmt.Errorf("échec écriture '%s' : %w", outPath, err)
		}
		outputPaths[i] = outPath

		mu.Lock()
		converted++
		current := converted
		mu.Unlock()

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

	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
