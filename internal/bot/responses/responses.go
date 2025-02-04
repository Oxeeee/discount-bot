package responses

import (
	"fmt"
	"log/slog"

	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	// SendTextMessage send custom message
	SendTextMessage(chatID int64, text string, role domain.UserRole)

	//SendErrorMessage sends type message ("❌ Ошибка: +errText")
	SendErrorMessage(chatID int64, errText string, role domain.UserRole)
	SendNotEnoughRightsMessage(chatID int64, role domain.UserRole)
	SendWelcomeMessage(chatID int64, role domain.UserRole)
	SendHelpMessage(chatID int64, role domain.UserRole)
	SendSuccessMessage(chatID int64, data any, role domain.UserRole)
}

type responder struct {
	bot *tgbotapi.BotAPI
	log *slog.Logger
}

func NewResponder(bot *tgbotapi.BotAPI, log *slog.Logger) Sender {
	return &responder{bot: bot, log: log}
}

func (r *responder) SendTextMessage(chatID int64, text string, role domain.UserRole) {
	const op = "SendTextMessage"
	r.log.With(slog.String("op", op))

	keyboard := GetKeyboard(role)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	_, err := r.bot.Send(msg)
	if err != nil {
		r.log.Error("error while sending message", "error", err)
	}

}

func (r *responder) SendErrorMessage(chatID int64, errText string, role domain.UserRole) {
	r.SendTextMessage(chatID, "❌ Ошибка: "+errText, role)
}

func (r *responder) SendNotEnoughRightsMessage(chatID int64, role domain.UserRole) {
	r.SendTextMessage(chatID, "⛔ У вас нет прав для выполнения этой команды", role)
}

func (r *responder) SendWelcomeMessage(chatID int64, role domain.UserRole) {
	r.SendTextMessage(chatID, "Hello, welcome to the bot! Type /help to see available commands.", role)
}

func (r *responder) SendHelpMessage(chatID int64, role domain.UserRole) {
	text := `/start - Начать работу с ботом
/adduser {username} - Добавить нового пользователя
/removeuser {username} - Удалить пользователя
/broadcast {message} - Отправить сообщение всем пользователям
/help - Показать список команд`
	r.SendTextMessage(chatID, text, role)
}

func (r *responder) SendSuccessMessage(chatID int64, data any, role domain.UserRole) {
	r.SendTextMessage(chatID, fmt.Sprintf("✅ %v", data), role)
}
