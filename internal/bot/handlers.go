package bot

import (
	"log/slog"

	"github.com/Oxeeee/discont-bot/internal/bot/responses"
	"github.com/Oxeeee/discont-bot/internal/domain"
	"github.com/Oxeeee/discont-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type CommandHandler func(b *BotHandler, message *tgbotapi.Message)
type ButtonHandler func(b *BotHandler, message *tgbotapi.Message)

type BotHandler struct {
	bot                *tgbotapi.BotAPI
	commandHandlers    map[string]CommandHandler
	buttonHandlers     map[string]ButtonHandler
	forceReplyHandlers map[int]ButtonHandler
	pendingRoleChange  map[int64]uint
	service            services.UserService
	log                *slog.Logger
	sender             responses.Sender
}

func NewBotHandler(bot *tgbotapi.BotAPI, service services.UserService, log *slog.Logger, sender responses.Sender) *BotHandler {
	bh := &BotHandler{
		bot:                bot,
		service:            service,
		log:                log,
		sender:             sender,
		forceReplyHandlers: make(map[int]ButtonHandler),
		pendingRoleChange:  make(map[int64]uint),
	}
	bh.commandHandlers = map[string]CommandHandler{
		"start": bh.handleStartCommand,
	}
	bh.buttonHandlers = map[string]ButtonHandler{
		"Проверить код":             StaffMiddleware(bh.handleCheckCodeButton),
		"Получить код":              UserMiddleware(bh.handleGetDiscountCodeButton),
		"Показать список скидок":    UserMiddleware(bh.handleShowDiscountsButton),
		"Управление пользователями": AdminMiddleware(bh.handleManageUsersButton),
		"🔁 Изменить список скидок":  AdminMiddleware(bh.handleEditDiscountsButton),
		"✅ Добавить пользователя":   AdminMiddleware(bh.handleAddUserButton),
		"❌ Удалить пользователя":    AdminMiddleware(bh.handleDeleteUserButton),
		"🔁 Поменять роль":           AdminMiddleware(bh.handleChangeRoleButton),
		"📋 Список пользователей":    AdminMiddleware(bh.handleUserListButton),
		"⬅️ Назад":                  UserMiddleware(bh.handleBackButton),
		"Пользователь":              AdminMiddleware(bh.handleRoleSelection),
		"Сотрудник":                 AdminMiddleware(bh.handleRoleSelection),
		"Администратор":             AdminMiddleware(bh.handleRoleSelection),
	}
	return bh
}

func (b *BotHandler) HandleMessage(message *tgbotapi.Message) {
	b.log.Info("Новое сообщение", slog.Int64("chat_id", message.Chat.ID), slog.Int64("user_id", message.From.ID), slog.String("username", message.From.UserName), slog.String("text", message.Text))
	if message.ReplyToMessage != nil {
		if handler, exists := b.forceReplyHandlers[message.ReplyToMessage.MessageID]; exists {
			handler(b, message)
			delete(b.forceReplyHandlers, message.ReplyToMessage.MessageID)
			return
		}
	}
	if handler, exists := b.buttonHandlers[message.Text]; exists {
		handler(b, message)
		return
	}
	if handler, exists := b.commandHandlers[message.Command()]; exists {
		handler(b, message)
	} else {
		b.log.Warn("Неизвестная команда", slog.String("command", message.Command()))
	}
}

func (b *BotHandler) handleStartCommand(_ *BotHandler, message *tgbotapi.Message) {
	exists, _, err := b.service.GetUserByID(uint(message.From.ID))
	if exists {
		b.updateUserKeyboard(message.Chat.ID, message.From.ID)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		b.log.Error("error while finding user by id", "error", err)
	}
	newUser := domain.User{
		ID:        uint(message.From.ID),
		Username:  message.From.UserName,
		FirstName: message.From.FirstName,
		LastName:  message.From.LastName,
		ChatID:    message.Chat.ID,
		Role:      string(domain.UserRoleUser),
		Whitelist: false,
	}
	err = b.service.RegisterUser(&newUser)
	if err != nil {
		b.log.Error("error while register user", "error", err)
		userRole, err := b.service.GetUserRole(uint(message.From.ID))
		if err != nil {
			b.log.Error("error while get user role", "error", err)
		}
		b.sender.SendErrorMessage(message.Chat.ID, err.Error(), domain.UserRole(userRole))
		return
	}
	b.updateUserKeyboard(message.Chat.ID, message.From.ID)
}

func (b *BotHandler) handleBackButton(_ *BotHandler, message *tgbotapi.Message) {
	b.updateUserKeyboard(message.Chat.ID, message.From.ID)
}
