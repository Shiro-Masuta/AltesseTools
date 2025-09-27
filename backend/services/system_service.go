package services

import (
	"context"
	"time"

	"Altesse_Tools_V1.0/backend/internal/system"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SystemService struct {
	ctx context.Context
}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (s *SystemService) Startup(ctx context.Context) {
	s.ctx = ctx

	// Goroutine qui envoie des infos toutes les 2 secondes
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				info, err := system.GetSystemMinimalInfo()
				if err != nil {
					continue
				}
				// Envoi d'un événement au frontend
				runtime.EventsEmit(ctx, "system:update", info)
			}
		}
	}()
}

// Méthode callable depuis le front pour un refresh manuel
func (s *SystemService) GetNow() (*system.SystemMinimalInfo, error) {
	return system.GetSystemMinimalInfo()
}
