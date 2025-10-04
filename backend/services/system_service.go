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
}

// --- Méthodes exposées au frontend ---

// Lancer la boucle d’updates (appelée depuis le frontend)
func (s *SystemService) Start() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				info, err := system.GetSystemMinimalInfo()
				if err != nil {
					continue
				}
				runtime.EventsEmit(s.ctx, "system:update", info)
			}
		}
	}()
}

// Snapshot immédiat
func (s *SystemService) GetNow() (*system.SystemMinimalInfo, error) {
	return system.GetSystemMinimalInfo()
}
