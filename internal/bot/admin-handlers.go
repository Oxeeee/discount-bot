package bot

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Oxeeee/discont-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) handleEditDiscountsButton(_ *BotHandler, message *tgbotapi.Message) {
	csvContent, err := b.service.GetDiscountList()
	if err != nil {
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка при получении списка пользователей", domain.UserRoleAdmin)
		b.log.Error("error while get discount list", "error", err)
		return
	}

	csvBytes := []byte(csvContent)

	doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
		Name:  "places.csv",
		Bytes: csvBytes,
	})

	doc.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}
	doc.Caption = "Отправьте новый CSV файл со списком скидок"
	sentMsg, err := b.bot.Send(doc)
	if err != nil {
		b.log.Error("error while sending csv doc", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка отправки CSV файла", domain.UserRoleAdmin)
		return
	}
	b.forceReplyHandlers[sentMsg.MessageID] = b.handleCSVFileReply
}

func (b *BotHandler) handleManageUsersButton(_ *BotHandler, message *tgbotapi.Message) {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("✅ Добавить пользователя"), tgbotapi.NewKeyboardButton("🔁 Поменять роль")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("❌ Удалить пользователя"), tgbotapi.NewKeyboardButton("📋 Список пользователей"), tgbotapi.NewKeyboardButton("⬅️ Назад")),
	)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = kb
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("Error while sending message", slog.Any("error", err))
	}
}

func (b *BotHandler) handleAddUserButton(_ *BotHandler, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите никнейм пользователя Пример: (@petrushin_leonid):")
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
	sentMsg, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("error while sending user nickname", slog.Any("error", err))
		return
	}
	b.forceReplyHandlers[sentMsg.MessageID] = b.handleAddUser
}

func (b *BotHandler) handleAddUser(_ *BotHandler, message *tgbotapi.Message) {
	nickname := message.Text
	nickname = strings.TrimPrefix(nickname, "@")
	userRole, err := b.service.GetUserRole(uint(message.From.ID))
	if err != nil {
		b.sender.SendErrorMessage(message.Chat.ID, err.Error(), domain.UserRole(userRole))
		b.log.Error("error while get user role", "error", err)
		return
	}
	if err := b.service.ManageWhitelist(nickname, "add"); err != nil {
		b.log.Error("error while whitlistening user", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, err.Error(), domain.UserRole(userRole))
		return
	}
	b.sender.SendSuccessMessage(message.Chat.ID, fmt.Sprintf("Пользователь %s добавлен", nickname), domain.UserRole(userRole))
}

func (b *BotHandler) handleChangeRoleButton(_ *BotHandler, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите никнейм пользователя (например: @petrushin_leonid):")
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
	sentMsg, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("Ошибка отправки запроса на ввод никнейма", slog.Any("error", err))
		return
	}
	b.forceReplyHandlers[sentMsg.MessageID] = b.handleChangeRoleAskForNickname
}

func (b *BotHandler) handleChangeRoleAskForNickname(_ *BotHandler, message *tgbotapi.Message) {
	nickname := strings.TrimSpace(message.Text)
	nickname = strings.TrimPrefix(nickname, "@")
	if nickname == "" {
		b.sender.SendErrorMessage(message.Chat.ID, "empty username", domain.UserRoleAdmin)
		return
	}
	exists, targetUser, err := b.service.GetUserByUsername(nickname)
	if err != nil || !exists {
		b.sender.SendErrorMessage(message.Chat.ID, fmt.Sprintf("Пользователь @%s не найден", nickname), domain.UserRoleAdmin)
		return
	}
	b.pendingRoleChange[message.Chat.ID] = targetUser.ID
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пользователь"),
			tgbotapi.NewKeyboardButton("Сотрудник"),
			tgbotapi.NewKeyboardButton("Администратор"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
		),
	)
	promptText := fmt.Sprintf("Выберите новую роль для пользователя @%s:", nickname)
	msg2 := tgbotapi.NewMessage(message.Chat.ID, promptText)
	msg2.ReplyMarkup = kb
	_, err = b.bot.Send(msg2)
	if err != nil {
		b.log.Error("Ошибка отправки сообщения с выбором роли", slog.Any("error", err))
	}
}

func (b *BotHandler) handleRoleSelection(_ *BotHandler, message *tgbotapi.Message) {
	targetUserID, exists := b.pendingRoleChange[message.Chat.ID]
	if !exists {
		return
	}
	selectedRole := message.Text
	if selectedRole == "Сотрудник" {
		selectedRole = "staff"
	} else if selectedRole == "Администратор" {
		selectedRole = "admin"
	} else if selectedRole == "Пользователь" {
		selectedRole = "user"
	}
	if err := b.service.ChangeRole(targetUserID, selectedRole); err != nil {
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка при изменении роли", domain.UserRoleAdmin)
	} else {
		b.sender.SendSuccessMessage(message.Chat.ID, "Роль успешно изменена", domain.UserRoleAdmin)
	}
	delete(b.pendingRoleChange, message.Chat.ID)
}

func (b *BotHandler) handleRoleSelectionBack(_ *BotHandler, message *tgbotapi.Message) {
	if _, exists := b.pendingRoleChange[message.Chat.ID]; exists {
		delete(b.pendingRoleChange, message.Chat.ID)
	}
	b.updateUserKeyboard(message.Chat.ID, message.From.ID)
}

func (b *BotHandler) handleDeleteUserButton(_ *BotHandler, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите никнейм пользователя\n\t\tПример: (@petrushin_leonid):")
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
	sentMsg, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("error while sending request for user nickname", slog.Any("error", err))
		return
	}
	b.forceReplyHandlers[sentMsg.MessageID] = b.handleDeleteUser
}

func (b *BotHandler) handleDeleteUser(_ *BotHandler, message *tgbotapi.Message) {
	nickname := message.Text
	nickname = strings.TrimPrefix(nickname, "@")
	if err := b.service.ManageWhitelist(nickname, "delete"); err != nil {
		b.log.Error("error while remove user from whitelist", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, err.Error(), domain.UserRoleAdmin)
		return
	}
	b.sender.SendSuccessMessage(message.Chat.ID, fmt.Sprintf("Пользователь %s удален", nickname), domain.UserRoleAdmin)
}

func (b *BotHandler) handleUserListButton(_ *BotHandler, message *tgbotapi.Message) {
	listText, err := b.service.Userlist()
	userRole, err := b.service.GetUserRole(uint(message.From.ID))
	if err != nil {
		b.log.Error("error while get user role", "error", err)
	}
	if err != nil {
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка при получении списка пользователей", domain.UserRole(userRole))
		b.updateUserKeyboard(message.Chat.ID, message.From.ID)
		b.log.Error("error while getting userlist", "error", err)
		return
	}
	b.sender.SendTextMessage(message.Chat.ID, listText, domain.UserRole(userRole))
}

func (b *BotHandler) handleCSVFileReply(_ *BotHandler, message *tgbotapi.Message) {
	if message.Document == nil {
		b.sender.SendErrorMessage(message.Chat.ID, "Пожалуйста, отправьте CSV файл.", domain.UserRoleAdmin)
		return
	}
	if !strings.HasSuffix(strings.ToLower(message.Document.FileName), ".csv") {
		b.sender.SendErrorMessage(message.Chat.ID, "Неверный формат файла. Ожидается CSV.", domain.UserRoleAdmin)
		return
	}
	fileConfig := tgbotapi.FileConfig{FileID: message.Document.FileID}
	file, err := b.bot.GetFile(fileConfig)
	if err != nil {
		b.log.Error("Ошибка получения файла", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, "Не удалось получить файл", domain.UserRoleAdmin)
		return
	}
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.bot.Token, file.FilePath)
	resp, err := http.Get(fileURL)
	if err != nil {
		b.log.Error("Ошибка загрузки файла", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, "Не удалось загрузить файл", domain.UserRoleAdmin)
		return
	}
	defer resp.Body.Close()
	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		b.log.Error("Ошибка чтения файла", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, "Не удалось прочитать файл", domain.UserRoleAdmin)
		return
	}
	csvContent := string(contentBytes)
	err = b.service.SaveDiscountList(csvContent)
	if err != nil {
		b.log.Error("Ошибка сохранения CSV", "error", err)
		b.sender.SendErrorMessage(message.Chat.ID, "Ошибка сохранения CSV", domain.UserRoleAdmin)
		return
	}
	b.sender.SendSuccessMessage(message.Chat.ID, "CSV файл успешно сохранен", domain.UserRoleAdmin)
}
