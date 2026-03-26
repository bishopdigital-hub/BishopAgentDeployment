package main

import (
	"embed"
	"log"
	"net/http"
	"fmt"

	"github.com/bishopcore/remote-agent/internal/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/bishopcore/remote-agent/internal/relay"
	"io/fs"
	"context"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/icon.ico
var iconData []byte

var defaultRole = "unified"
var defaultDiscordToken = "MTQ4NTAxOTc0Mzc1MTc2NjA5Nw.GX_2yC.C9wqyeNIGF7zdysZhiSUCIKTDKn2_UvvLL-Y5c"
var defaultDiscordChannel = "1485020950457094398"

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("⚠️ Config failed to load: %v, using defaults", err)
		cfg = &config.Config{
			NexusURL: "ws://41.193.213.128:9500/api/remote/agent",
			Role:     defaultRole,
			WebPort:  9600,
		}
	}

	// Create an instance of the app structure
	app := NewApp(cfg)

	// Start browser-accessible UI server
	go func() {
		subDist, err := fs.Sub(assets, "frontend/dist")
		if err != nil {
			log.Printf("⚠️ Asset sub-dir pivot failed: %v", err)
			return
		}
		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.FS(subDist)))
		log.Printf("🌐 Browser access enabled at http://127.0.0.1:%d", cfg.WebPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.WebPort), mux); err != nil {
			log.Printf("⚠️ HTTP Server failed: %v", err)
		}
	}()

	// Initialize Discord Relay (v1.5)
	token := cfg.DiscordToken
	if token == "" {
		token = defaultDiscordToken
	}
	channel := cfg.DiscordChannelID
	if channel == "" {
		channel = defaultDiscordChannel
	}

	if token != "" {
		dr := relay.NewDiscordRelay(token, channel, cfg.AgentID, cfg.AgentName, cfg.Role)
		go dr.Start(context.Background())
	}

	// Setup Tray (Conditional Compilation)
	go SetupTray(app)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Bishop Remote Agent",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 2, G: 6, B: 23, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		HideWindowOnClose: true,
		StartHidden:       true,
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
