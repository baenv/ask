package reg

import (
	"ask/pkg/config"
	"ask/pkg/logger"
	"ask/pkg/models"
	"ask/pkg/repo"
	"ask/pkg/utils/encryptutils"
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/go-telegram/bot"
	telegramMod "github.com/go-telegram/bot/models"
)

type Telegram struct {
	config        config.Config
	repo          repo.Repository
	logger        logger.Logger
	lastHandlerID string
}

func NewTelegram(repo repo.Repository, config config.Config, logger logger.Logger) *Telegram {
	return &Telegram{
		repo:   repo,
		config: config,
		logger: logger,
	}
}

func (t *Telegram) Handle(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	b.UnregisterHandler(t.lastHandlerID)

	if !isPermissable(b, update) {
		t.logger.Warn("Unauthorized access attempt")
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "You don't have permission to perform this action.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send permission denied message")
		}
		return
	}

	if update.Message != nil && update.Message.Text != "" {
		t.handleCommand(ctx, b, update)
	} else if update.CallbackQuery != nil {
		t.handleCallbackQuery(ctx, b, update)
	}
}

func (t *Telegram) handleCommand(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	command := update.Message.Text

	switch command {
	case "/reg":
		t.handleUserRegistration(ctx, b, update)
	case "/reg server":
		t.handleServerRegistration(ctx, b, update)
	default:
		return
	}
}

func (t *Telegram) handleCallbackQuery(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	data := update.CallbackQuery.Data
	parts := strings.SplitN(data, ":", 2)

	if len(parts) != 2 {
		return
	}

	action, configID := parts[0], parts[1]

	switch action {
	case "reg_setup_user":
		t.handleUserSetup(ctx, b, update, configID)
	case "reg_setup_server":
		t.handleServerSetup(ctx, b, update, configID)
	case "reg_remove_user":
		t.handleUserRemoval(ctx, b, update, configID)
	case "reg_remove_server":
		t.handleServerRemoval(ctx, b, update, configID)
	}
}

func (t *Telegram) handleUserSetup(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	config, err := t.repo.UserConfig().GetByID(id)
	if err != nil {
		t.logger.Error(err, "Failed to get user config")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "An error occurred. Please try again.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send error message")
		}
		return
	}

	if config.Command == "" {
		t.fillUserCommand(ctx, b, update, id)
		return
	}

	if config.EndpointURL == "" {
		t.fillUserEndpointURL(ctx, b, update, id)
		return
	}

	if config.APIKey == "" {
		t.fillUserApiKey(ctx, b, update, id)
		return
	}

	if config.Description == "" {
		t.fillUserDescription(ctx, b, update, id)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Setup is already complete.",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send setup complete message")
	}
}

func (t *Telegram) fillUserCommand(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		command := update.Message.Text
		err := t.repo.UserConfig().SaveCommand(id, command)
		if err != nil {
			t.logger.Error(err, "Failed to save user command")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "Failed to save command. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "Command saved.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send command saved message")
		}

		b.UnregisterHandler(t.lastHandlerID)
		t.fillUserEndpointURL(ctx, b, update, id)
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Please enter the command for this configuration:",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send command prompt")
	}
}

func (t *Telegram) fillUserEndpointURL(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		endpointURL := update.Message.Text
		if !isValidURL(endpointURL) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "Invalid endpoint URL. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send invalid URL message")
			}
			return
		}

		err := t.repo.UserConfig().SaveEndpointURL(id, endpointURL)
		if err != nil {
			t.logger.Error(err, "Failed to save user endpoint URL")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "Failed to save endpoint URL. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "Endpoint URL saved.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send URL saved message")
		}

		b.UnregisterHandler(t.lastHandlerID)
		t.fillUserApiKey(ctx, b, update, id)
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Please enter the endpoint URL:",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send endpoint URL prompt")
	}
}

func (t *Telegram) fillUserApiKey(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		apiKey := update.Message.Text

		// Encrypt the API key before saving
		encryptionKey, err := encryptutils.NewEncryptionKey(t.config.EncryptionKey)
		if err != nil {
			t.logger.Error(err, "Failed to create encryption key")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "An error occurred. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		encryptedAPIKey, err := encryptutils.EncryptAPIKey(encryptionKey, apiKey)
		if err != nil {
			t.logger.Error(err, "Failed to encrypt API key")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "An error occurred. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		err = t.repo.UserConfig().SaveAPIKey(id, encryptedAPIKey)
		if err != nil {
			t.logger.Error(err, "Failed to save user API key")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "Failed to save API key. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "API key saved.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send API key saved message")
		}

		b.UnregisterHandler(t.lastHandlerID)
		t.fillUserDescription(ctx, b, update, id)

		// delete update message
		_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
		if err != nil {
			t.logger.Error(err, "Failed to delete message")
		}
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Please enter your API key:",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send API key prompt")
	}
}

func (t *Telegram) fillUserDescription(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		description := update.Message.Text
		err := t.repo.UserConfig().SaveDescription(id, description)
		if err != nil {
			t.logger.Error(err, "Failed to save user description")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "Failed to save description. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "Description saved. Setup complete.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send setup complete message")
		}

		b.UnregisterHandler(t.lastHandlerID)
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Please enter a description for this configuration:",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send description prompt")
	}
}

func (t *Telegram) handleServerSetup(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	config, err := t.repo.ServerConfig().GetByID(id)
	if err != nil {
		t.logger.Error(err, "Failed to get server config")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "An error occurred. Please try again.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send error message")
		}
		return
	}

	if config.Command == "" {
		t.fillServerCommand(ctx, b, update, id)
		return
	}

	if config.EndpointURL == "" {
		t.fillServerApiUrl(ctx, b, update, id)
		return
	}

	if config.APIKey == "" {
		t.fillServerApiKey(ctx, b, update, id)
		return
	}

	if config.Description == "" {
		t.fillServerDescription(ctx, b, update, id)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Setup is already complete.",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send setup complete message")
	}
}

func (t *Telegram) fillServerCommand(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		command := update.Message.Text
		err := t.repo.ServerConfig().SaveCommand(id, command)
		if err != nil {
			t.logger.Error(err, "Failed to save server command")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "Failed to save command. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "Command saved.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send command saved message")
		}

		b.UnregisterHandler(t.lastHandlerID)
		t.fillServerApiUrl(ctx, b, update, id)
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Please enter the command for this server configuration:",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send command prompt")
	}
}

func (t *Telegram) fillServerDescription(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		description := update.Message.Text
		err := t.repo.ServerConfig().SaveDescription(id, description)
		if err != nil {
			t.logger.Error(err, "Failed to save config description")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: getUserID(update),
				Text:   "Failed to save config description. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: getUserID(update),
			Text:   "Description saved. Setup complete.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send setup complete message")
		}

		b.UnregisterHandler(t.lastHandlerID)
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: getUserID(update),
		Text:   "Please enter a description for this configuration:",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send description prompt")
	}
}

func (t *Telegram) fillServerApiUrl(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		endpointURL := update.Message.Text
		if !isValidURL(endpointURL) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Invalid endpoint URL. Please try again.",
			})
			return
		}

		err := t.repo.ServerConfig().SaveEndpointURL(id, endpointURL)
		if err != nil {
			t.logger.Error(err, "Failed to save server endpoint URL")
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Failed to save endpoint URL. Please try again.",
			})
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Endpoint URL saved.",
		})

		b.UnregisterHandler(t.lastHandlerID)
		t.fillServerApiKey(ctx, b, update, id)
	})

	var chatID int64
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.From.ID
	} else if update.Message != nil {
		chatID = update.Message.Chat.ID
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Please enter the endpoint URL:",
	})
}

func (t *Telegram) fillServerApiKey(ctx context.Context, b *bot.Bot, update *telegramMod.Update, id string) {
	b.UnregisterHandler(t.lastHandlerID)
	t.lastHandlerID = b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
		apiKey := update.Message.Text

		// Encrypt the API key before saving
		encryptionKey, err := encryptutils.NewEncryptionKey(t.config.EncryptionKey)
		if err != nil {
			t.logger.Error(err, "Failed to create encryption key")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "An error occurred. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		encryptedAPIKey, err := encryptutils.EncryptAPIKey(encryptionKey, apiKey)
		if err != nil {
			t.logger.Error(err, "Failed to encrypt API key")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "An error occurred. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		err = t.repo.ServerConfig().SaveAPIKey(id, encryptedAPIKey)
		if err != nil {
			t.logger.Error(err, "Failed to save server API key")
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Failed to save API key. Please try again.",
			})
			if err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "API key saved.",
		})
		if err != nil {
			t.logger.Error(err, "Failed to send API key saved message")
		}

		b.UnregisterHandler(t.lastHandlerID)
		t.fillServerDescription(ctx, b, update, id)

		// delete update message
		_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
		if err != nil {
			t.logger.Error(err, "Failed to delete message")
		}
	})

	var chatID int64
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.From.ID
	} else if update.Message != nil {
		chatID = update.Message.Chat.ID
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Please enter the secret token:",
	})
	if err != nil {
		t.logger.Error(err, "Failed to send secret token prompt")
	}
}

func (t *Telegram) handleUserRemoval(ctx context.Context, b *bot.Bot, update *telegramMod.Update, userID string) {
	err := t.repo.UserConfig().RemoveByID(userID)
	if err != nil {
		t.logger.Error(err, "Failed to remove user config")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Failed to remove user configuration. Please try again.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:   "User configuration removed successfully.",
	})
}

func (t *Telegram) handleServerRemoval(ctx context.Context, b *bot.Bot, update *telegramMod.Update, serverID string) {
	err := t.repo.ServerConfig().RemoveByID(serverID)
	if err != nil {
		t.logger.Error(err, "Failed to remove server config")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Failed to remove server configuration. Please try again.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:   "Server configuration removed successfully.",
	})
}

func (t *Telegram) handleUserRegistration(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	chatID := update.Message.Chat.ID
	userID := fmt.Sprintf("%d", update.Message.From.ID)

	pendingConfigs, err := t.repo.UserConfig().ListPendingByUserID(userID)
	if err != nil {
		t.logger.Error(err, "Failed to get pending user configs")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Failed to retrieve pending configurations. Please try again.",
		})
		return
	}

	if len(pendingConfigs) > 0 {
		// Prepare inline keyboard with pending configurations
		var keyboard [][]tgbotapi.InlineKeyboardButton
		setupCallbackData := fmt.Sprintf("reg_setup_user:%d", pendingConfigs[0].ID)
		removeCallbackData := fmt.Sprintf("reg_remove_user:%d", pendingConfigs[0].ID)

		setupButton := tgbotapi.NewInlineKeyboardButtonData("Setup", setupCallbackData)
		removeButton := tgbotapi.NewInlineKeyboardButtonData("❌", removeCallbackData)

		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(setupButton, removeButton))

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboard...)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        "You have a pending configuration. Please complete it or remove it:",
			ReplyMarkup: inlineKeyboard,
		})
		return
	}

	user := models.User{
		UserID:   userID,
		Username: update.Message.From.Username,
		Platform: models.PlatformTelegram,
		UserAgentConfigs: []models.UserAgentConfig{
			{
				APIKey:      "",
				EndpointURL: "",
			},
		},
	}

	user, err = t.repo.User().Create(user)
	if err != nil {
		t.logger.Error(err, "Failed to create user")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Failed to register user. Please try again.",
		})
		return
	}

	// Create inline keyboard for user configuration
	setupCallbackData := fmt.Sprintf("reg_setup_user:%d", user.UserAgentConfigs[0].ID)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Setup Configuration", setupCallbackData),
		),
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        "User registered successfully. Please complete your configuration by tapping the button below:",
		ReplyMarkup: keyboard,
	})
}

func (t *Telegram) handleServerRegistration(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	chatID := update.Message.Chat.ID
	userID := fmt.Sprintf("%d", update.Message.From.ID)
	isDM := chatID >= 0

	var err error
	if !isDM {
		// If in a group chat, create a new server configuration
		user := models.User{
			UserID:   userID,
			Username: update.Message.From.Username,
			Platform: models.PlatformTelegram,
			Servers: []models.Server{
				{
					ServerID:   fmt.Sprintf("%d", chatID),
					Platform:   models.PlatformTelegram,
					ServerName: update.Message.Chat.Title,
					ServerAdminConfig: []models.ServerAdminConfig{
						{APIKey: "", EndpointURL: ""},
					},
				},
			},
		}
		_, err = t.repo.User().Create(user)
		if err != nil {
			t.logger.Error(err, "Failed to create server")
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Failed to register server. Please try again.",
			})
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Server configuration has been initiated successfully for this chat.",
		})
	}

	// Get pending configurations
	pendingConfigs, err := t.repo.ServerConfig().ListPendingByUserID(userID)
	if err != nil {
		t.logger.Error(err, "Failed to get pending servers")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Failed to retrieve pending servers. Please try again.",
		})
		return
	}

	if len(pendingConfigs) == 0 && isDM {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "You have no pending server configurations. If you want to setup a new server, please move to the server chat and type the command `/reg server`.",
		})
		return
	}
	// Prepare inline keyboard with pending configurations
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, config := range pendingConfigs {
		setupCallbackData := fmt.Sprintf("reg_setup_server:%d", config.ID)
		removeCallbackData := fmt.Sprintf("reg_remove_server:%d", config.ID)

		setupButton := tgbotapi.NewInlineKeyboardButtonData(config.Server.ServerName, setupCallbackData)
		removeButton := tgbotapi.NewInlineKeyboardButtonData("❌", removeCallbackData)

		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(setupButton, removeButton))
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	// Send setup message with inline keyboard
	recipientID := chatID
	if !isDM {
		recipientID = update.Message.From.ID
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		Text:        "Please select a server to complete its configuration or remove it:",
		ChatID:      recipientID,
		ReplyMarkup: inlineKeyboard,
	})
}
