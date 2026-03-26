//go:build !windows

package main

import (
	"log"
)

func SetupTray(app *App) {
	log.Println("ℹ️ System Tray disabled on this platform. Background persistence active via Wails runtime.")
}
