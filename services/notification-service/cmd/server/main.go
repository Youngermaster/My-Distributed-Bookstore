package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/youngermaster/distributed-bookstore/notification-service/internal/config"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/logger"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/notification"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/queue"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	cfg.Print()

	baseLogger := logger.New(cfg.LogLevel)

	history := notification.NewHistory(200)

	dispatcher := notification.NewDispatcher(
		baseLogger.With().Str("component", "dispatcher").Logger(),
		[]notification.Sender{
			notification.NewLogSender(notification.ChannelEmail, baseLogger.With().Str("channel", "email").Logger()),
			notification.NewLogSender(notification.ChannelSMS, baseLogger.With().Str("channel", "sms").Logger()),
			notification.NewLogSender(notification.ChannelPush, baseLogger.With().Str("channel", "push").Logger()),
		},
		history,
	)

	consumer, err := queue.NewConsumer(cfg, dispatcher, baseLogger.With().Str("component", "consumer").Logger())
	if err != nil {
		baseLogger.Fatal().Err(err).Msg("failed to initialize rabbitmq consumer")
	}
	defer consumer.Close()

	httpServer := server.New(cfg, history, baseLogger.With().Str("component", "http").Logger())

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(rootCtx); err != nil && rootCtx.Err() == nil {
			baseLogger.Error().Err(err).Msg("consumer stopped unexpectedly")
			cancel()
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%s", cfg.HTTPPort)
		baseLogger.Info().Str("addr", addr).Msg("starting http server")
		if err := httpServer.Listen(addr); err != nil && rootCtx.Err() == nil {
			baseLogger.Fatal().Err(err).Msg("http server stopped unexpectedly")
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		baseLogger.Info().Str("signal", sig.String()).Msg("shutdown signal received")
	case <-rootCtx.Done():
		baseLogger.Info().Msg("context cancelled, shutting down")
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		baseLogger.Error().Err(err).Msg("failed to shutdown http server cleanly")
	}

	baseLogger.Info().Msg("notification service stopped")
}
