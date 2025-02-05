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

	isValid, username, codeInfo, err := b.service.VerifyCode(code)
	if err != nil {
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка проверки кода", domain.UserRoleStaff)
		b.updateUserKeyboard(message.Chat.ID, message.From.ID)
		return
	}
	if isValid {
		b.sender.SendTextMessage(message.Chat.ID, fmt.Sprintf("✅ Код действителен!\n Пользователь: %v\n Заканчивается в %v", username, codeInfo.ExpDate), domain.UserRoleStaff)
	} else {
		b.sender.SendTextMessage(message.Chat.ID, "❌ Код недействителен!", domain.UserRoleStaff)
	}
	b.updateUserKeyboard(message.Chat.ID, message.From.ID)
}
