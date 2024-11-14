package start

import (
	"context"
	"strings"
	"sum/pkg/logger"
	"sum/pkg/repo"

	"github.com/go-telegram/bot"
	telegramMod "github.com/go-telegram/bot/models"
)

type Telegram struct {
	repo   repo.Repository
	logger logger.Logger
}

func NewTelegram(repo repo.Repository, logger logger.Logger) *Telegram {
	return &Telegram{
		repo:   repo,
		logger: logger,
	}
}

func (t *Telegram) Handle(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	command := strings.ToLower(update.Message.Text)
	if command == "/start" || command == "/help" {
		t.sendHelpMessage(ctx, b, update)
	}
}

func (t *Telegram) sendHelpMessage(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	helpText := `👋 Welcome to ask Bot! Here are the available commands:

🚀 /start or /help
   Show this help message

📝 /reg
   Register a new user or server configuration
   • Set up a new user configuration
   • /reg server - Set up a new server configuration

📋 /ls
   List your commands or server configurations
   • List your personal commands
   • /ls server - List server configurations (in private chat) or server commands (in group chat)

🤖 /ai <subcommand> <message>
   Execute an AI command
   • Format: /ai <subcommand> <message>
   • The subcommand should match one of your configured commands

📌 Example: /ai summarize Please summarize this text for me.

-------------------------------------------

❓ Need more help? Feel free to ask!`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	})

	if err != nil {
		t.logger.Error(err, "Failed to send help message")
	}
}
