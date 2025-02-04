package responses

import (
	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetKeyboard(role domain.UserRole) tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton

	switch role {
	case domain.UserRoleUser:
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Получить код")},
			{tgbotapi.NewKeyboardButton("Показать список скидок")},
		}
	case domain.UserRoleStaff:
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Проверить код")},
		}
	case domain.UserRoleAdmin:
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Управление пользователями")},
			{tgbotapi.NewKeyboardButton("🔁 Изменить список скидок")},
		}
	}

	return tgbotapi.NewReplyKeyboard(buttons...)
}
