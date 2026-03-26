package app

import (
	"log"
)

type AgentRole string

const (
	RoleServer      AgentRole = "server"
	RoleWorkstation AgentRole = "workstation"
)

type Config struct {
	Role       AgentRole `json:"role"`
	NexusURL   string    `json:"nexus_url"`
	AuthToken  string    `json:"auth_token"`
	DeviceID   string    `json:"device_id"`
	TrelloBoard string   `json:"trello_board"`
}

type App struct {
	Config *Config
	// Add relays and services here
}

func NewApp(cfg *Config) *App {
	return &App{
		Config: cfg,
	}
}

func (a *App) Start() error {
	log.Printf("[ ⟐ BISHOP_CORE ] Initializing Remote Agent in %s mode...", a.Config.Role)
	// Initialize relay and start background services
	return nil
}
