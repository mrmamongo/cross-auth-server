package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrmamongo/go-auth-tg/pkg/logger"
)

type CommandHandler func(*tgbotapi.BotAPI, *tgbotapi.Message) error

type MessageHandler func(*tgbotapi.BotAPI, *tgbotapi.Message) error

type TelegramBotDispatcher struct {
	bot             *tgbotapi.BotAPI
	l               logger.Interface
	commandHandlers map[string]CommandHandler
	messageHandler  MessageHandler
}

func NewTelegramBotDispatcher(bot *tgbotapi.BotAPI, l logger.Interface) *TelegramBotDispatcher {
	return &TelegramBotDispatcher{
		bot:             bot,
		l:               l,
		commandHandlers: make(map[string]CommandHandler),
		messageHandler: func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) error {
			return nil
		},
	}
}

func (c *TelegramBotDispatcher) Update(update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	if update.Message.IsCommand() {
		command := update.Message.Command()
		handler, ok := c.commandHandlers[command]
		if !ok {
			return fmt.Errorf("unknown command: %s", command)
		}
		return handler(c.bot, update.Message)
	}
	return c.messageHandler(c.bot, update.Message)
}

func (c *TelegramBotDispatcher) RegisterCommandHandler(command string, handler CommandHandler) {
	c.commandHandlers[command] = handler
}

func (c *TelegramBotDispatcher) RegisterMessageHandler(handler MessageHandler) {
	c.messageHandler = handler
}
