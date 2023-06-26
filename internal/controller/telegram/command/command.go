package command

import (
	"context"
	"fmt"
	botapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
		dispatcher.RegisterCommandHandler("/start", c.StartCommand)
		dispatcher.RegisterCommandHandler("/change", c.ChangeUsernameCommand)
	}
}

func (c *UserCommands) StartCommand(ctx context.Context, bot *botapi.Bot, message *models.Message) {
	c.l.Info("Received /start command")
	var user *entity.User
	var err error
	if user, err = c.u.GetByTelegram(context.Background(), message.From.Username); user == nil {
		_, err = bot.SendMessage(ctx, &botapi.SendMessageParams{
			ChatID: message.Chat.ID,
			Text:   fmt.Sprintf("Sorry, you are not registered. Please, register at the website")},
		)
		if err != nil {
			c.l.Error(err)
			return
		}
		return
	}
	_, err = bot.SendMessage(ctx, &botapi.SendMessageParams{
		ChatID: message.Chat.ID,
		Text:   fmt.Sprintf("Hello, %s!", user.Username),
	})
	if err != nil {
		c.l.Error(err)
		return
	}
	c.l.Info("User %s has started!", user.Username)
	if err != nil {
		c.l.Error(err)
	}
}

func (c *UserCommands) ChangeUsernameCommand(ctx context.Context, bot *botapi.Bot, message *models.Message) {
	c.l.Info("Received /change command")
	var user *entity.User
	var err error

	if user, err = c.u.GetByTelegram(context.Background(), message.Chat.Username); err != nil {
		c.l.Error(err)
		return
	}

	if user == nil {
		_, err = bot.SendMessage(ctx, &botapi.SendMessageParams{
			ChatID: message.Chat.ID,
			Text:   fmt.Sprintf("Sorry, you are not registered. Please, register at the website")},
		)
		if err != nil {
			c.l.Error(err)
			return
		}
	}

	user.Username = strings.Replace(message.Text, "/change ", "", 1)
	if err = c.u.Update(context.Background(), user); err != nil {
		c.l.Error(err)
		return
	}

	_, err = bot.SendMessage(ctx, &botapi.SendMessageParams{
		ChatID: message.Chat.ID,
		Text:   fmt.Sprintf("Hello, %s!", user.Username),
	})
	if err != nil {
		c.l.Error(err)
	}
}
