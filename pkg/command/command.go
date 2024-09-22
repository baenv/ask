// Package command provides functionality for handling commands across different platforms.
package command

import (
	"ask/pkg/adapter"
	"ask/pkg/config"
	"ask/pkg/logger"
	"ask/pkg/repo"

	"gorm.io/gorm"

	"github.com/bwmarrin/discordgo"
	"github.com/go-telegram/bot"
)

// Command represents a structure that holds command handlers for different platforms.
type Command struct {
	Discord  ICommand // Discord command handler
	Telegram ICommand // Telegram command handler
}

// New creates a new Command instance with initialized Discord and Telegram handlers.
// It takes the application configuration, Discord session, and Telegram bot as parameters.
func New(cfg config.Config, d *discordgo.Session, t *bot.Bot, logger logger.Logger, db *gorm.DB) Command {
	a := adapter.New(cfg)
	repo := repo.NewRepository(db)

	return Command{
		Discord:  NewDiscord(repo, d, a, logger),
		Telegram: NewTelegram(repo, t, cfg, a, logger),
	}
}
