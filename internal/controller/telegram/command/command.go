package command

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrmamongo/go-auth-tg/internal/entity"
	"github.com/mrmamongo/go-auth-tg/internal/usecase"
	"github.com/mrmamongo/go-auth-tg/pkg/logger"
	"github.com/mrmamongo/go-auth-tg/pkg/tgbot"
	"strings"
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
		dispatcher.RegisterCommandHandler("change", c.ChangeUsernameCommand)
	}
}

func (c *UserCommands) StartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	c.l.Info("Received /start command")
	var user *entity.User
	var err error
	if user, err = c.u.GetByTelegram(context.Background(), message.Chat.UserName); err != nil {
		return err
	}
	_, err = bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Hello, %s!", user.Username)))
	if err != nil {
		return err
	}
	c.l.Info("User %s has started!", user.Username)
	return nil
}

func (c *UserCommands) ChangeUsernameCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	c.l.Info("Received /change command")
	var user *entity.User
	var err error

	if user, err = c.u.GetByTelegram(context.Background(), message.Chat.UserName); err != nil {
		return err
	}

	if user == nil {
		_, err = bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Sorry, you are not registered. Please, register at the website")))
		if err != nil {
			return err
		}
	}

	user.Username = strings.Replace(message.Text, "/change ", "", 1)
	if err = c.u.Update(context.Background(), user); err != nil {
		return err
	}

	_, err = bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Hello, %s!", user.Username)))
	return err
}
