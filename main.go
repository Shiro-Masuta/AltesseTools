package main

import (
	"context"
	"embed"

	service "Altesse_Tools_V1.0/backend/services"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	converter := service.NewConverterService()
	system := service.NewSystemService()
	stats := service.NewStatsService()
	rename := service.NewRenameService()
	duplicate := service.NewDuplicateService()
	err := wails.Run(&options.App{
		Title:     "Altesse_Tools_V1.0",
		Width:     1250,
		Height:    800,
		MaxWidth:  1250,
		MaxHeight: 800,
		MinWidth:  800,
		MinHeight: 500,
		Frameless: true,
		LogLevel:  logger.INFO,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			converter.SetContext(ctx)
			system.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			converter,
			system,
			stats,
			rename,
			duplicate,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
