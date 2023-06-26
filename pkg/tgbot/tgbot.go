package tgbot

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/mrmamongo/go-auth-tg/config"
	"net/http"
)

type TelegramBot struct {
	bot        *bot.Bot
	dispatcher *TelegramBotDispatcher
	notify     chan error
	handler    http.HandlerFunc
}

func NewTelegramBot(ctx context.Context, api *bot.Bot, dispatcher *TelegramBotDispatcher) *TelegramBot {
	b := &TelegramBot{bot: api, notify: make(chan error), dispatcher: dispatcher}
	b.handler = b.bot.WebhookHandler()
	b.start(ctx)

	return b
}

func (t *TelegramBot) start(ctx context.Context) {
	go t.bot.StartWebhook(ctx)
}

func (t *TelegramBot) Shutdown() {
	//t.bot.StopReceivingUpdates()
}

func (t *TelegramBot) Notify() chan error {
	return t.notify
}

func (t *TelegramBot) webhookHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t.handler(
			c.Writer, c.Request,
		)
	}
}

func (t *TelegramBot) RegisterRouter(handler *gin.Engine, config *config.Config) {
	handler.POST("/api/v1/bot"+config.Token, t.webhookHandler())
}
