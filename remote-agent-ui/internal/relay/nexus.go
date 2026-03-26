package relay

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type NexusRelay struct {
	URL  string
	conn *websocket.Conn
}

func NewNexusRelay(url string) *NexusRelay {
	return &NexusRelay{URL: url}
}

func (n *NexusRelay) Connect(ctx context.Context) error {
	if strings.HasPrefix(n.URL, "http") {
		// If HTTP URL provided, assume ws is on /api/ws (legacy) or just use http for telemetry
		return nil 
	}
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(n.URL, nil)
	if err != nil {
		return err
	}
	n.conn = conn
	log.Printf("📡 Connected to Nexus at %s", n.URL)
	return nil
}

func (n *NexusRelay) SendTelemetry(payload interface{}) error {
	// 1. WebSocket attempt (Legacy Sync)
	if n.conn != nil {
		n.conn.WriteJSON(payload)
	}

	// 2. HTTP POST attempt (V1.5+ CommandSync)
	endpoint := strings.Replace(n.URL, "ws://", "http://", 1)
	endpoint = strings.Replace(endpoint, "wss://", "https://", 1)
	endpoint = strings.TrimSuffix(endpoint, "/api/remote/agent")
	endpoint += "/api/command/sync"

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️ HTTP Sync failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️ HTTP Sync returned %d", resp.StatusCode)
	}

	return nil
}
