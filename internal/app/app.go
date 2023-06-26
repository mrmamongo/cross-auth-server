// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	botapi "github.com/go-telegram/bot"
	"github.com/mrmamongo/go-auth-tg/internal/controller/telegram/command"
	"github.com/mrmamongo/go-auth-tg/pkg/tgbot"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/mrmamongo/go-auth-tg/config"
	v1 "github.com/mrmamongo/go-auth-tg/internal/controller/http/v1"
	"github.com/mrmamongo/go-auth-tg/internal/usecase"
	"github.com/mrmamongo/go-auth-tg/internal/usecase/repo"
	"github.com/mrmamongo/go-auth-tg/pkg/httpserver"
	"github.com/mrmamongo/go-auth-tg/pkg/logger"
	"github.com/mrmamongo/go-auth-tg/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	userUseCase := usecase.NewUserUseCase(
		repo.NewUserRepo(pg),
	)

	// Telegram Bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	b, err := botapi.New(cfg.Token)
	_, err = b.DeleteWebhook(ctx, &botapi.DeleteWebhookParams{DropPendingUpdates: true})

	_, err = b.SetWebhook(ctx, &botapi.SetWebhookParams{
		URL: "https://45c8-2a00-1370-81a8-8fe-e5d0-297b-f366-4365.ngrok-free.app/api/v1/bot5895361951:AAEytOyi3g_sNEigRmFHQNfUI2bXdMWXT1U",
	})
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	dispatcher := tgbot.NewTelegramBotDispatcher(b, l)
	command.NewUserCommands(dispatcher, l, userUseCase)

	bot := tgbot.NewTelegramBot(ctx, b, dispatcher)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, userUseCase)
	bot.RegisterRouter(handler, cfg)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-bot.Notify():
		l.Error(fmt.Errorf("app - Run - bot.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	bot.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - bot.Shutdown: %w", err))
	}
}
