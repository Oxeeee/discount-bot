package bot

import (
	"log/slog"

	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AdminMiddleware(handler CommandHandler) func(b *BotHandler, message *tgbotapi.Message) {
	return func(b *BotHandler, message *tgbotapi.Message) {
		const op = "admin middleware"
		b.log.With(slog.String("op", op))

		userID := message.From.ID
		userRole, err := b.service.GetUserRole(uint(message.From.ID))
		if err != nil {
			b.log.Error("error while get user role", "error", err)
			return
		}

		if userRole != string(domain.UserRoleAdmin) {
			b.log.Warn("user is not an admin", "user_id", userID)
			b.sender.SendNotEnoughRightsMessage(message.Chat.ID, domain.UserRole(userRole))
			return
		}
		handler(b, message)
	}
}

func StaffMiddleware(handler CommandHandler) func(b *BotHandler, message *tgbotapi.Message) {
	return func(b *BotHandler, message *tgbotapi.Message) {
		const op = "staff middleware"
		b.log.With(slog.String("op", op))

		userID := message.From.ID
		userRole, err := b.service.GetUserRole(uint(message.From.ID))
		if err != nil {
			b.log.Error("error while get user role", "error", err)
			return
		}

		if userRole != string(domain.UserRoleStaff) {
			b.log.Warn("user is not an staff", "user_id", userID)
			b.sender.SendNotEnoughRightsMessage(message.Chat.ID, domain.UserRole(userRole))
			return
		}
		handler(b, message)
	}
}

func UserMiddleware(handler CommandHandler) func(b *BotHandler, message *tgbotapi.Message) {
	return func(b *BotHandler, message *tgbotapi.Message) {
		const op = "user middleware"
		b.log.With(slog.String("op", op))

		userID := message.From.ID
		isWhitelisted, err := b.service.CheckWhitelist(uint(userID))
		if err != nil {
			b.log.Error("Ошибка проверки whitelist", "error", err)
			return
		}

		if !isWhitelisted {
			b.log.Warn("Пользователь не в whitelist", "user_id", userID)
			b.sender.SendNotEnoughRightsMessage(message.Chat.ID, domain.UserRoleUser)
			return
		}

		handler(b, message)
	}
}
