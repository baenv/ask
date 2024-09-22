// Package listener provides functionality for handling different messaging platforms.
package listener

import (
	"ask/pkg/command"
	"context"

	"github.com/go-telegram/bot"
)

// telegram represents a Telegram listener instance
type telegram struct {
	bot     *bot.Bot         // Telegram bot instance
	command command.ICommand // Command handler for Telegram
}

// NewTelegram initiates a Telegram listener instance
func NewTelegram(b *bot.Bot, c command.ICommand) IListener {
	return &telegram{
		bot:     b,
		command: c,
	}
}

// Start begins the Telegram bot in a separate goroutine
func (t telegram) Start() error {
	go t.bot.Start(context.Background())
	return nil
}

// End stops the Telegram bot (currently a no-op)
func (t telegram) End() error {
	return nil
}

// Register registers the sum command for Telegram
func (t *telegram) Register() {
	// t.command.RegisterReg()
	// t.command.RegisterLs()
	// t.command.RegisterAi()
	// t.command.RegisterStart()
	t.command.RegisterAsk()
}
