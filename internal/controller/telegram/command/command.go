package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrmamongo/go-auth-tg/internal/usecase"
	"github.com/mrmamongo/go-auth-tg/pkg/logger"
	"github.com/mrmamongo/go-auth-tg/pkg/tgbot"
)

type UserCommands struct {
	u usecase.User
	l logger.Interface
}

func NewUserCommands(dispatcher *tgbot.TelegramBotDispatcher, l logger.Interface, u usecase.User) {
	c := &UserCommands{
		u: u,
		l: l,
	}

	{
		dispatcher.RegisterCommandHandler("start", c.StartCommand)
		dispatcher.RegisterCommandHandler("register", c.RegisterCommand)
	}
}

func (c *UserCommands) StartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	c.l.Info("Received /start command")
	return nil
}

func (c *UserCommands) RegisterCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	c.l.Info("Received /register command")
	return nil
}
