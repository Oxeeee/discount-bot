
package bot

import (
	"log/slog"

	"github.com/Oxeeee/discont-bot/internal/bot/responses"
	"github.com/Oxeeee/discont-bot/internal/config"
	"github.com/Oxeeee/discont-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	log     *slog.Logger
	botAPI  *tgbotapi.BotAPI
	handler *BotHandler
	cfg     *config.Config
}

func New(log *slog.Logger, cfg *config.Config, service services.UserService) *App {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		panic(err)
	}

	botAPI.Debug = false

	sender := responses.NewResponder(botAPI, log)
	handler := NewBotHandler(botAPI, service, log, sender)

	return &App{
		log:     log,
		botAPI:  botAPI,
		handler: handler,
	}
}

func (a *App) MustRun() {
	const op = "bot.run"
	a.log.With(slog.String("op", op))

	a.log.Info("Telegram bot is running", slog.String("username", a.botAPI.Self.UserName))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := a.botAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			a.handler.HandleMessage(update.Message)
		}
	}
}
