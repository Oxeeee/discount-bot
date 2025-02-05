package bot

import (
	"fmt"
	"log/slog"

	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) handleCheckCodeButton(_ *BotHandler, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите код для проверки:")
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
	sentMsg, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("Ошибка отправки запроса на ввод кода", slog.Any("error", err))
		return
	}
	b.forceReplyHandlers[sentMsg.MessageID] = b.handleCodeVerification
}

func (b *BotHandler) handleCodeVerification(_ *BotHandler, message *tgbotapi.Message) {
	code := message.Text

	isValid, user, err := b.service.VerifyCode(code, uint(message.From.ID))
	if err != nil {
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка проверки кода", domain.UserRoleStaff)
		b.updateUserKeyboard(message.Chat.ID, message.From.ID)
		return
	}

	if isValid {
		if user.LastName == "" {
			user.LastName = "Не указано"
		}
		user.Username = "@" + user.Username
		if user.CodesUsed >= 1 {
			user.CodesUsed = user.CodesUsed - 1
		}

		b.sender.SendTextMessage(message.Chat.ID, fmt.Sprintf("✅ Код действителен!\nПользователь: %v,\nИмя: %v,\nФамилия: %v,\nКоличество использованных кодов: %v", user.Username, user.LastName, user.FirstName, user.CodesUsed), domain.UserRoleStaff)
	} else {
		b.sender.SendTextMessage(message.Chat.ID, "❌ Код недействителен!", domain.UserRoleStaff)
	}
	b.updateUserKeyboard(message.Chat.ID, message.From.ID)
}
