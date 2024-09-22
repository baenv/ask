package ai

import (
	"ask/pkg/adapter"
	"ask/pkg/config"
	"ask/pkg/logger"
	"ask/pkg/models"
	"ask/pkg/repo"
	"ask/pkg/utils/encryptutils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	telegramMod "github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

type Telegram struct {
	repo    repo.Repository
	logger  logger.Logger
	adapter adapter.IAdapter
	config  config.Config
}

func NewTelegram(repo repo.Repository, config config.Config, adapter adapter.IAdapter, logger logger.Logger) *Telegram {
	return &Telegram{
		repo:    repo,
		logger:  logger,
		adapter: adapter,
		config:  config,
	}
}

func (t *Telegram) Handle(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	parts := strings.Fields(update.Message.Text)
	if len(parts) < 3 {
		return
	}

	subcommand := parts[1]
	message := strings.Join(parts[2:], " ")

	userID := fmt.Sprintf("%d", update.Message.From.ID)
	user, _ := t.repo.User().GetByPlatformID(userID, "telegram")

	chatID := fmt.Sprintf("%d", update.Message.Chat.ID)
	server, _ := t.repo.Server().GetByPlatformID(chatID, "telegram")

	var config interface{}
	var err error
	if update.Message.Chat.Type == "private" {
		config, err = t.repo.UserConfig().GetByUserIDAndCommand(user.ID, subcommand)
	} else {
		config, err = t.repo.ServerConfig().GetByServerIDAndCommand(server.ID, subcommand)
		if err != nil {
			config, err = t.repo.UserConfig().GetByUserIDAndCommand(user.ID, subcommand)
		}
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   escapeSpecialChars("Command configuration not found. Try /ls or /ls server to check if the command is set up."),
			}); err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
		} else {
			t.logger.Error(err, "Failed to retrieve command config")
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   escapeSpecialChars("Failed to retrieve command configuration. Please try again."),
			}); err != nil {
				t.logger.Error(err, "Failed to send error message")
			}
		}
		return
	}

	if config == nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   escapeSpecialChars(fmt.Sprintf("Command '%s' not found.", subcommand)),
		}); err != nil {
			t.logger.Error(err, "Failed to send error message")
		}
		return
	}

	// Execute the command
	url := ""
	encryptedToken := ""
	switch c := config.(type) {
	case models.UserAgentConfig:
		url = c.EndpointURL
		encryptedToken = c.APIKey
	case models.ServerAdminConfig:
		url = c.EndpointURL
		encryptedToken = c.APIKey
	}

	// Decrypt the API key
	encryptionKey, err := encryptutils.NewEncryptionKey(t.config.EncryptionKey)
	if err != nil {
		t.logger.Error(err, "Failed to create encryption key")
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   escapeSpecialChars("An error occurred. Please try again."),
		}); err != nil {
			t.logger.Error(err, "Failed to send error message")
		}
		return
	}

	token, err := encryptutils.DecryptAPIKey(encryptionKey, encryptedToken)
	if err != nil {
		t.logger.Error(err, "Failed to decrypt API key")
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   escapeSpecialChars("An error occurred. Please try again."),
		}); err != nil {
			t.logger.Error(err, "Failed to send error message")
		}
		return
	}

	response, err := chat(t.adapter, message, url, token)
	if err != nil {
		t.logger.Error(err, "Error executing command")
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   escapeSpecialChars(fmt.Sprintf("Error executing command: %v", err)),
		}); err != nil {
			t.logger.Error(err, "Failed to send error message")
		}
		return
	}
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      response.Summary,
		ParseMode: telegramMod.ParseModeMarkdownV1,
	}); err != nil {
		t.logger.Error(err, "Failed to send message")
		// Attempt to send an error message to the user
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   escapeSpecialChars("An error occurred while sending the message. Please try again."),
		}); err != nil {
			t.logger.Error(err, "Failed to send error message")
		}
	}
}

// Sum represents the structure of a summarized article
type Sum struct {
	URL     string
	Title   string
	Summary string
}

func chat(a adapter.IAdapter, msg, url, token string) (*Sum, error) {
	data, err := a.Dify().Chat(msg, url, token)
	if err != nil {
		return nil, fmt.Errorf("failed to summarize the article: %w", err)
	}

	if data == "" {
		return nil, errors.New("empty response from LLM")
	}

	// Process the response
	lines := strings.Split(data, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	processedData := strings.Join(lines, "\n")

	return &Sum{
		URL:     url,
		Summary: processedData,
	}, nil
}

func escapeSpecialChars(s string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		if !strings.Contains(s, "\\"+char) {
			s = strings.ReplaceAll(s, char, "\\"+char)
		}
	}
	return s
}
