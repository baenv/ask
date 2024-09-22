// Package command provides functionality for handling commands across different platforms.
package command

import (
	"ask/pkg/adapter"
	"ask/pkg/command/ai"
	"ask/pkg/command/ask"
	"ask/pkg/command/ls"
	"ask/pkg/command/reg"
	"ask/pkg/command/start"
	"ask/pkg/config"
	"ask/pkg/logger"
	"ask/pkg/repo"

	"github.com/go-telegram/bot"
)

// telegram represents a Telegram command handler.
type telegram struct {
	bot   *bot.Bot
	reg   *reg.Telegram
	ls    *ls.Telegram
	ai    *ai.Telegram
	start *start.Telegram
	ask   *ask.Telegram
}

// NewTelegram creates a new Telegram command handler.
func NewTelegram(repo repo.Repository, t *bot.Bot, cfg config.Config, a adapter.IAdapter, logger logger.Logger) ICommand {
	return &telegram{
		bot:   t,
		reg:   reg.NewTelegram(repo, cfg, logger),
		ls:    ls.NewTelegram(repo, logger),
		ai:    ai.NewTelegram(repo, cfg, a, logger),
		start: start.NewTelegram(repo, logger),
		ask:   ask.NewTelegram(cfg, a, logger),
	}
}

// AddHandler adds the command handler to the Telegram bot.
// Currently, this method is empty and can be implemented as needed.
func (t *telegram) AddHandler() {}

// RegisterReg registers the reg command with the Telegram bot.
func (t *telegram) RegisterReg() {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/reg", bot.MatchTypePrefix, t.reg.Handle)
	t.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, "reg_", bot.MatchTypePrefix, t.reg.Handle)
}

// RegisterLs registers the ls command with the Telegram bot.
func (t *telegram) RegisterLs() {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/ls", bot.MatchTypePrefix, t.ls.Handle)
	t.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, "ls_", bot.MatchTypePrefix, t.ls.Handle)
}

// RegisterAI registers the ai command with the Telegram bot.
func (t *telegram) RegisterAi() {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/ai", bot.MatchTypePrefix, t.ai.Handle)
}

// RegisterStart registers the start command with the Telegram bot.
func (t *telegram) RegisterStart() {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, t.start.Handle)
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypePrefix, t.start.Handle)
}

// RegisterAsk registers the ask command with the Telegram bot.
func (t *telegram) RegisterAsk() {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/ask", bot.MatchTypePrefix, t.ask.Handle)
}
