package ask

import (
	"ask/pkg/adapter"
	"ask/pkg/config"
	"ask/pkg/logger"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	telegramMod "github.com/go-telegram/bot/models"
)

type Telegram struct {
	logger  logger.Logger
	adapter adapter.IAdapter
	config  config.Config
}

func NewTelegram(config config.Config, adapter adapter.IAdapter, logger logger.Logger) *Telegram {
	return &Telegram{
		logger:  logger,
		adapter: adapter,
		config:  config,
	}
}

// Handle executes the /ask command and sends the response to the user.
// It also logs any errors that occur while executing the command.
func (t *Telegram) Handle(ctx context.Context, b *bot.Bot, update *telegramMod.Update) {
	parts := strings.Fields(update.Message.Text)
	if len(parts) < 2 {
		return
	}

	message := strings.Join(parts[1:], " ")

	// Create a channel to signal when the response is ready
	responseChan := make(chan *Sum, 1)
	errChan := make(chan error, 1)

	// Start a goroutine to fetch the response
	go func() {
		response, err := chat(t.adapter, message, t.config.AgentURL, t.config.AgentToken)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	}()

	// Send a "thinking" message after 1 second if no response yet
	var thinkingMsg *telegramMod.Message
	select {
	case <-time.After(1 * time.Second):
		var err error
		thinkingMsg, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:              update.Message.Chat.ID,
			Text:                "🤔 Thinking...",
			ProtectContent:      true,
			DisableNotification: true,
			ReplyParameters: &telegramMod.ReplyParameters{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			},
		})
		if err != nil {
			t.logger.Error(err, "Failed to send thinking message")
		}
	case response := <-responseChan:
		sendResponse(ctx, b, update, response, t.logger)
		return
	case err := <-errChan:
		sendErrorMessage(ctx, b, update, err, t.logger)
		return
	}

	// Wait for the actual response
	select {
	case response := <-responseChan:
		// Delete the "thinking" message if it was sent
		if thinkingMsg != nil {
			if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: thinkingMsg.ID,
			}); err != nil {
				t.logger.Error(err, "Failed to delete thinking message")
			}
		}
		sendResponse(ctx, b, update, response, t.logger)
	case err := <-errChan:
		sendErrorMessage(ctx, b, update, err, t.logger)
	}
}

func sendResponse(ctx context.Context, b *bot.Bot, update *telegramMod.Update, response *Sum, logger logger.Logger) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:              update.Message.Chat.ID,
		Text:                response.Summary,
		ParseMode:           telegramMod.ParseModeMarkdownV1,
		ProtectContent:      true,
		DisableNotification: true,
		ReplyParameters: &telegramMod.ReplyParameters{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		},
	}); err != nil {
		logger.Error(err, "Failed to send message")
		sendErrorMessage(ctx, b, update, errors.New("An error occurred while sending the message. Please try again."), logger)
	}
}

func sendErrorMessage(ctx context.Context, b *bot.Bot, update *telegramMod.Update, err error, logger logger.Logger) {
	logger.Error(err, "Error executing command")
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:              update.Message.Chat.ID,
		Text:                fmt.Sprintf("Error executing command: %v", err),
		ProtectContent:      true,
		DisableNotification: true,
		ReplyParameters: &telegramMod.ReplyParameters{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		},
	}); err != nil {
		logger.Error(err, "Failed to send error message")
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

	return &Sum{
		URL:     url,
		Summary: data,
	}, nil
}
