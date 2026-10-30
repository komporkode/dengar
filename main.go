package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func init() {
	if os.Getenv("WAILS_DEV_ALLOWED_HOSTS") == "" {
		_ = os.Setenv("WAILS_DEV_ALLOWED_HOSTS", "wails.localhost,localhost,127.0.0.1")
	}
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Dengerin - YouTube Music Desktop",
		Width:     1200,
		Height:    800,
		Frameless: true,
		MinWidth:  360,
		MinHeight: 140,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 9, G: 10, B: 14, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
