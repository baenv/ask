package listener

import (
	"sum/pkg/command"

	"github.com/bwmarrin/discordgo"
)

// discord represents a Discord listener instance
type discord struct {
	session *discordgo.Session // Discord session
	command command.ICommand   // Command handler for Discord
}

// NewDiscord initiates a Discord listener instance
func NewDiscord(s *discordgo.Session, c command.ICommand) IListener {
	return &discord{
		session: s,
		command: c,
	}
}

// Start opens the Discord session
func (d discord) Start() error {
	return d.session.Open()
}

// End closes the Discord session
func (d discord) End() error {
	return d.session.Close()
}

// Register registers the sum command for Discord
func (d *discord) Register() {
	d.command.RegisterReg()
}
