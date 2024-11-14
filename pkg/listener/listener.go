package listener

import (
	"sum/pkg/command"
	"sum/pkg/config"
	"sum/pkg/logger"

	"gorm.io/gorm"

	"github.com/bwmarrin/discordgo"
	"github.com/go-telegram/bot"
)

// IListener defines the interface for a listener component.
// It provides methods to start and stop the listener.
type IListener interface {
	Start() error
	End() error
	Register()
}

// Listener is a struct that contains supported platform listener
type Listener struct {
	Discord  IListener
	Telegram IListener
}

// New creates an instance of Listener
func New(cfg config.Config, logger logger.Logger, d *discordgo.Session, t *bot.Bot, db *gorm.DB) Listener {
	command := command.New(cfg, d, t, logger, db)

	var discord, telegram IListener
	if d != nil {
		command.Discord.AddHandler()
		discord = NewDiscord(d, command.Discord)
	}
	if t != nil {
		telegram = NewTelegram(t, command.Telegram)
	}

	return Listener{
		Discord:  discord,
		Telegram: telegram,
	}
}
