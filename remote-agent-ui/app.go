package main

import (
	"context"
	"fmt"
	"time"
	"github.com/bishopcore/remote-agent/internal/config"
	"github.com/bishopcore/remote-agent/internal/relay"
	"github.com/bishopcore/remote-agent/internal/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	cfg     *config.Config
	nexus   *relay.NexusRelay
	scanner *service.ScanService
	shell   *service.ShellService
	audit   *service.AuditService
	trello  *service.TrelloService
	exploit *service.ExploitService
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Config) *App {
	return &App{
		cfg:     cfg,
		nexus:   relay.NewNexusRelay(cfg.NexusURL),
		scanner: service.NewScanService(),
		shell:   service.NewShellService(),
		audit:   service.NewAuditService(),
		trello:  service.NewTrelloService("REPLACE_ME", "REPLACE_ME"),
		exploit: service.NewExploitService("127.0.0.1", 55553, "msf", "afro"),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.nexus.Connect(ctx)
	go a.telemetryLoop(ctx)
}

func (a *App) telemetryLoop(ctx context.Context) {
	// Immediate first sync
	a.syncTelemetry()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.syncTelemetry()
		}
	}
}

func (a *App) syncTelemetry() {
	payload := map[string]interface{}{
		"command": "HEARTBEAT",
		"status":  "online",
		"source":  a.cfg.AgentID,
		"company": a.cfg.Company,
		"version": "1.5.1",
		"payload": map[string]interface{}{
			"cpu": 0.0,
			"ram": 0.0,
			"net": 0.0,
		},
	}
	a.nexus.SendTelemetry(payload)
}

// Login verifies the operator password
func (a *App) Login(password string) bool {
	if a.cfg.OperatorPassword == "" {
		// If no password set, allow login (for initial setup)
		return true
	}
	return password == a.cfg.OperatorPassword
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// RunScan triggers an Nmap scan
func (a *App) RunScan(target string) string {
	res, err := a.scanner.ExecuteScan(target)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return res
}

// RunAudit triggers a Lynis audit
func (a *App) RunAudit() string {
	res, err := a.audit.RunLynisAudit()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return res
}

// RunExploit triggers an autonomous exploit (Mock)
func (a *App) RunExploit(target string) string {
	// For demo: try to authenticate and list modules
	if err := a.exploit.Authenticate(); err != nil {
		return fmt.Sprintf("Exploit Auth Failed: %v", err)
	}
	return "MSF Authenticated. Exploitation sequence initiated (Simulated)."
}

// GetServiceStatus checks if a service is reachable (Mock for now)
func (a *App) GetServiceStatus(name string) string {
	// In a real scenario, we would check the port/process
	// For now, we simulate connectivity
	return "online"
}

// LaunchService opens a service URL in the default browser
func (a *App) LaunchService(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
