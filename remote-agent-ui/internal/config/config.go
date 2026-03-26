package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	AgentID          string `json:"agent_id"`
	AgentName        string `json:"agent_name"`
	Company          string `json:"company"`
	Version          string `json:"version"`
	Role             string `json:"role"`
	NexusURL         string `json:"nexus_url"`
	WebPort          int    `json:"web_port"`
	DiscordToken     string `json:"discord_token"`
	DiscordChannelID string `json:"discord_channel_id"`
	OperatorPassword string `json:"operator_password"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
