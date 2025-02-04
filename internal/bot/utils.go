package bot

import (
	"log/slog"

	"github.com/Oxeeee/discont-bot/internal/bot/responses"
	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) updateUserKeyboard(chatID int64, userID int64) {
	role, err := b.service.GetUserRole(uint(userID))
	if err != nil {
		b.log.Error("Ошибка получения роли пользователя", slog.Any("error", err))
		return
	}

	keyboard := responses.GetKeyboard(domain.UserRole(role))

	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
	msg.ReplyMarkup = keyboard

	_, err = b.bot.Send(msg)
	if err != nil {
		b.log.Error("Ошибка обновления клавиатуры", slog.Any("error", err))
	}
}
