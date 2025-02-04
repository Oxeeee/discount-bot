package bot

import (
	"fmt"

	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) handleGetDiscountCodeButton(_ *BotHandler, message *tgbotapi.Message) {
	code, err := b.service.GenerateCode(uint(message.From.ID))
	if err != nil {
		b.log.Error("error while generating discount code", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка при генерации кода, обратитесь к тех. администратору: @petrushin_leonid", domain.UserRoleUser)
		return
	}

	b.sender.SendTextMessage(message.Chat.ID, fmt.Sprintf("Ваш скидочный код %v, он действителен 15 минут", code), domain.UserRoleUser)
}

func (b *BotHandler) handleShowDiscountsButton(_ *BotHandler, message *tgbotapi.Message) {
	userRole, err := b.service.GetUserRole(uint(message.From.ID))
	if err != nil {
		b.log.Error("error while get user role", "error", err)
	}
	b.sender.SendTextMessage(message.Chat.ID, "Список текущих скидок:\n- 10% на кофе\n- 15% на одежду", domain.UserRole(userRole))
}
