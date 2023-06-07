// Package app configures and runs application.
package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrmamongo/go-auth-tg/internal/controller/telegram/command"
	"github.com/mrmamongo/go-auth-tg/pkg/tgbot"
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
	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, userUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Telegram Bot
	b, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - tgbotapi.NewBotAPI: %w", err))
	}
	dispatcher := tgbot.NewTelegramBotDispatcher(b, l)
	command.NewUserCommands(dispatcher, l, userUseCase)

	bot := tgbot.NewTelegramBot(b, dispatcher)

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
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}
}
