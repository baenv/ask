package ls

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sum/pkg/logger"
	"sum/pkg/models"
	"sum/pkg/repo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	if update.Message != nil && update.Message.Text != "" {
		t.handleMessage(ctx, b, update)
	} else if update.CallbackQuery != nil {
		t.handleCallbackQuery(ctx, b, update)
	}
}

func (t *Telegram) handleMessage(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	command := update.Message.Text
	userID := getUserID(update)

	switch command {
	case "/ls":
		t.listUserCommands(ctx, b, update, userID)
	case "/ls server":
		t.listServers(ctx, b, update, userID)
	}
}

func (t *Telegram) listUserCommands(ctx context.Context, b *bot.Bot, update *telegramMod.Update, userID int64) {
	user, err := t.repo.User().GetByPlatformID(fmt.Sprintf("%d", userID), string(models.PlatformTelegram))
	if err != nil {
		t.logger.Error(err, "Failed to retrieve user")
		t.sendErrorMessage(ctx, b, update, "Failed to retrieve user information\\. Please try again\\.")
		return
	}

	commands, err := t.repo.UserConfig().ListByUserID(user.ID)
	if err != nil {
		t.logger.Error(err, "Failed to retrieve user commands")
		t.sendErrorMessage(ctx, b, update, "Failed to retrieve user commands\\. Please try again\\.")
		return
	}

	if len(commands) == 0 {
		t.sendMessage(ctx, b, update, "You don't have any commands set up\\. Please use /reg to set up a command\\.")
		return
	}

	t.displayCommands(ctx, b, update, commands, "user")
}

func (t *Telegram) listServers(ctx context.Context, b *bot.Bot, update *telegramMod.Update, userID int64) {
	if update.Message.Chat.Type == "private" {
		user, err := t.repo.User().GetByPlatformID(fmt.Sprintf("%d", userID), string(models.PlatformTelegram))
		if err != nil {
			t.logger.Error(err, "Failed to retrieve user")
			t.sendErrorMessage(ctx, b, update, "Failed to retrieve user information\\. Please try again\\.")
			return
		}

		t.listUserServers(ctx, b, update, user.ID)
	} else {
		server, err := t.repo.Server().GetByPlatformID(fmt.Sprintf("%d", getChatID(update)), string(models.PlatformTelegram))
		if err != nil {
			t.logger.Error(err, "Failed to retrieve server")
			t.sendErrorMessage(ctx, b, update, "Failed to retrieve server information\\. Please try again\\.")
			return
		}

		t.listServerCommands(ctx, b, update, server.ID)
	}
}

func (t *Telegram) listUserServers(ctx context.Context, b *bot.Bot, update *telegramMod.Update, userID int64) {
	servers, err := t.repo.Server().ListByUserID(userID)
	if err != nil {
		t.logger.Error(err, "Failed to retrieve user servers")
		t.sendErrorMessage(ctx, b, update, "Failed to retrieve user servers\\. Please try again\\.")
		return
	}

	if len(servers) == 0 {
		t.sendMessage(ctx, b, update, "You don't have any servers registered\\. Please go to a server and use /reg server to register a server\\.")
		return
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, server := range servers {
		callbackData := fmt.Sprintf("ls_server:%d", server.ID)
		button := tgbotapi.NewInlineKeyboardButtonData(server.ServerName, callbackData)
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(button))
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	t.sendMessage(ctx, b, update, "Select a server to view its commands:", inlineKeyboard)
}

func (t *Telegram) listServerCommands(ctx context.Context, b *bot.Bot, update *telegramMod.Update, serverID int64) {
	commands, err := t.repo.ServerConfig().ListByServerID(serverID)
	if err != nil {
		t.logger.Error(err, "Failed to retrieve server commands")
		t.sendErrorMessage(ctx, b, update, "Failed to retrieve server commands\\. Please try again\\.")
		return
	}

	if len(commands) == 0 {
		t.sendMessage(ctx, b, update, "No commands found for this server\\. Please use /reg server to set up commands\\.")
		return
	}

	hasPermission := isPermissable(b, update)
	t.displayCommands(ctx, b, update, commands, "server", hasPermission)
}

func (t *Telegram) handleCallbackQuery(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	data := update.CallbackQuery.Data
	parts := strings.SplitN(data, ":", 2)

	if len(parts) != 2 {
		return
	}

	action, id := parts[0], parts[1]

	switch action {
	case "ls_server":
		var idInt int64
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return
		}

		t.listServerCommandsCallback(ctx, b, update, idInt)
	case "ls_remove_command":
		t.removeCommand(ctx, b, update, id)
	}
}

func (t *Telegram) listServerCommandsCallback(ctx context.Context, b *bot.Bot, update *telegramMod.Update, serverID int64) {
	commands, err := t.repo.ServerConfig().ListByServerID(serverID)
	if err != nil {
		t.logger.Error(err, "Failed to retrieve server commands")
		t.sendErrorMessage(ctx, b, update, "Failed to retrieve server commands\\. Please try again\\.")
		return
	}

	if len(commands) == 0 {
		t.sendMessage(ctx, b, update, "No commands found for this server\\. Please use /reg server in the appropriate server to set up commands\\.")
		return
	}

	hasPermission := isPermissable(b, update)
	t.displayCommands(ctx, b, update, commands, "server", hasPermission)
}

func (t *Telegram) displayCommands(ctx context.Context, b *bot.Bot, update *telegramMod.Update, commands interface{}, commandType string, hasPermission ...bool) {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var messageText string

	messageText += "📋 List of Commands\n\n"
	messageText += "Here are your current commands:\n\n"

	switch commandType {
	case "user":
		userCommands := commands.([]models.UserAgentConfig)
		for i, command := range userCommands {
			messageText += formatCommandInfo(command.Command, command.Description)
			if i < len(userCommands)-1 {
				messageText += "\n\\-\\-\\-\n\n"
			}
			removeCallbackData := fmt.Sprintf("ls_remove_command:user_%d", command.ID)
			removeButton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🗑️ Remove \"%s\"", command.Command), removeCallbackData)
			keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(removeButton))
		}
	case "server":
		serverCommands := commands.([]models.ServerAdminConfig)
		showRemoveButtons := len(hasPermission) > 0 && hasPermission[0]
		for i, command := range serverCommands {
			messageText += formatCommandInfo(command.Command, command.Description)
			if i < len(serverCommands)-1 {
				messageText += "\n\\-\\-\\-\n\n"
			}
			if showRemoveButtons {
				removeCallbackData := fmt.Sprintf("ls_remove_command:server_%d", command.ID)
				removeButton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🗑️ Remove \"%s\"", command.Command), removeCallbackData)
				keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(removeButton))
			}
		}
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	t.sendMessage(ctx, b, update, messageText, inlineKeyboard)
}

func (t *Telegram) removeCommand(ctx context.Context, b *bot.Bot, update *telegramMod.Update, commandID string) {
	if !isPermissable(b, update) {
		t.sendErrorMessage(ctx, b, update, "You don't have permission to remove this command\\.")
		return
	}

	parts := strings.SplitN(commandID, "_", 2)
	if len(parts) != 2 {
		return
	}

	commandType, id := parts[0], parts[1]

	var err error
	switch commandType {
	case "user":
		err = t.repo.UserConfig().RemoveByID(id)
	case "server":
		err = t.repo.ServerConfig().RemoveByID(id)
	}

	if err != nil {
		t.logger.Error(err, "Failed to remove command")
		t.sendErrorMessage(ctx, b, update, "Failed to remove command\\. Please try again\\.")
		return
	}

	t.sendMessage(ctx, b, update, "Command removed successfully\\.")
}

func formatCommandInfo(command, description string) string {
	return fmt.Sprintf("🤖 *Command:* `%s`\n📝 *Description:* %s", command, strings.ReplaceAll(description, "-", "\\-"))
}

func (t Telegram) sendMessage(ctx context.Context, b *bot.Bot, update *telegramMod.Update, text string, replyMarkup ...interface{}) {
	chatID := getChatID(update)
	params := &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: telegramMod.ParseModeMarkdown,
	}
	if len(replyMarkup) > 0 {
		params.ReplyMarkup = replyMarkup[0]
	}
	_, err := b.SendMessage(ctx, params)
	if err != nil {
		t.logger.Error(err, "Failed to send message")
	}
}

func (t Telegram) sendErrorMessage(ctx context.Context, b *bot.Bot, update *telegramMod.Update, text string) {
	t.sendMessage(ctx, b, update, text)
}
