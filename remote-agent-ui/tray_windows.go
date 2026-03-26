//go:build windows

package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func SetupTray(app *App) {
	systray.Run(func() {
		systray.SetIcon(iconData)
		systray.SetTitle("Bishop Remote")
		systray.SetTooltip("Bishop Unified Remote Agent")

		mPortal := systray.AddMenuItem("Open Portal", "Open agent interface in browser")
		mUpdate := systray.AddMenuItem("Update (v1.3)", "Check for version 1.3 updates")
		systray.AddSeparator()
		mShow := systray.AddMenuItem("Show Context", "Restore command window")
		mHide := systray.AddMenuItem("Minimize to Tray", "Hide window to background")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Terminate", "Shutdown agent service")

		for {
			select {
			case <-mPortal.ClickedCh:
				runtime.BrowserOpenURL(app.ctx, "http://localhost:9600")
			case <-mUpdate.ClickedCh:
				runtime.BrowserOpenURL(app.ctx, "https://dev.bishopdigital.co.za/bishop-agent-installer-v1.3.exe")
			case <-mShow.ClickedCh:
				runtime.WindowShow(app.ctx)
			case <-mHide.ClickedCh:
				runtime.WindowHide(app.ctx)
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}, func() {
		log.Println("Systray exiting")
	})
}
