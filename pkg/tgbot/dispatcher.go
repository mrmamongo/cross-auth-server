package tgbot

import (
	"context"
	"fmt"
	botapi "github.com/go-telegram/bot"
	models "github.com/go-telegram/bot/models"
	"github.com/mrmamongo/go-auth-tg/pkg/logger"
	"strings"
)

type CommandHandler func(ctx context.Context, bot *botapi.Bot, update *models.Message)

type MessageHandler func(ctx context.Context, bot *botapi.Bot, update *models.Message)

type TelegramBotDispatcher struct {
	l               logger.Interface
	commandHandlers map[string]CommandHandler
	messageHandler  MessageHandler
}

func NewTelegramBotDispatcher(bot *botapi.Bot, l logger.Interface) *TelegramBotDispatcher {
	dispatcher := TelegramBotDispatcher{
		l:               l,
		commandHandlers: make(map[string]CommandHandler),
		messageHandler: func(ctx context.Context, bot *botapi.Bot, update *models.Message) {
			return
		},
	}
	bot.RegisterHandler(botapi.HandlerTypeMessageText, "", botapi.MatchTypeContains, dispatcher.Update)
	return &dispatcher
}

func (c *TelegramBotDispatcher) Update(ctx context.Context, bot *botapi.Bot, update *models.Update) {
	if update.Message.Text == "" {
		return
	}

	if update.Message.Text[0] != '/' {
		c.messageHandler(ctx, bot, update.Message)
	}
	command := strings.Split(update.Message.Text, " ")[0]
	handler, ok := c.commandHandlers[command]
	if !ok {
		c.l.Error(fmt.Errorf("unknown command: %s", command))
		return
	}
	handler(ctx, bot, update.Message)
}

func (c *TelegramBotDispatcher) RegisterCommandHandler(command string, handler CommandHandler) {
	c.commandHandlers[command] = handler
}

func (c *TelegramBotDispatcher) RegisterMessageHandler(handler MessageHandler) {
	c.messageHandler = handler
}
