package ls

import (
	"context"
	"net/url"

	"github.com/go-telegram/bot"
	telegramMod "github.com/go-telegram/bot/models"
)

func isValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func isPermissable(b *bot.Bot, update *telegramMod.Update) bool {
	var fromID, chatID int64
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message.Message == nil {
			return false
		}

		fromID = update.CallbackQuery.From.ID
		chatID = update.CallbackQuery.Message.Message.Chat.ID
	} else {
		fromID = update.Message.From.ID
		chatID = update.Message.Chat.ID
	}

	if fromID == chatID {
		return true
	}

	member, err := b.GetChatMember(context.Background(), &bot.GetChatMemberParams{
		ChatID: chatID,
		UserID: fromID,
	})
	if err != nil {
		return false
	}

	return member.Type == telegramMod.ChatMemberTypeAdministrator || member.Type == telegramMod.ChatMemberTypeOwner
}

func getUserID(update *telegramMod.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return 0
}

func getChatID(update *telegramMod.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil {
		return update.CallbackQuery.Message.Message.Chat.ID
	}
	return 0
}
