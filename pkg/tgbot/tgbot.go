package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot        *tgbotapi.BotAPI
	dispatcher *TelegramBotDispatcher
	notify     chan error
}

func NewTelegramBot(api *tgbotapi.BotAPI, dispatcher *TelegramBotDispatcher) *TelegramBot {
	b := &TelegramBot{bot: api, notify: make(chan error), dispatcher: dispatcher}
	b.start()

	return b
}

func (t *TelegramBot) start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			err := t.dispatcher.Update(update)
			if err != nil {
				_, err = t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("tgbot - TelegramBot - Dispatch: %e", err)))
				if err != nil {
					return
				}
			}
		}
	}()
}

func (t *TelegramBot) Shutdown() {
	t.bot.StopReceivingUpdates()
}

func (t *TelegramBot) Notify() chan error {
	return t.notify
}
