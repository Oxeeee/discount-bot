package bot

import (
	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) handleGetDiscountCodeButton(_ *BotHandler, message *tgbotapi.Message) {
	userRole, err := b.service.GetUserRole(uint(message.From.ID))
	if err != nil {
		b.log.Error("error while get user role", "error", err)
	}
	b.sender.SendTextMessage(message.Chat.ID, "Ваш персональный скидочный код: ABC123", domain.UserRole(userRole))
}

func (b *BotHandler) handleShowDiscountsButton(_ *BotHandler, message *tgbotapi.Message) {
	userRole, err := b.service.GetUserRole(uint(message.From.ID))
	if err != nil {
		b.log.Error("error while get user role", "error", err)
	}
	b.sender.SendTextMessage(message.Chat.ID, "Список текущих скидок:\n- 10% на кофе\n- 15% на одежду", domain.UserRole(userRole))
}
