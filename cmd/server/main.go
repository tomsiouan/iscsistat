package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iscsistat/internal/config"

	"go.uber.org/zap"
)

// shutdownHandler manages the graceful shutdown of the HTTP server.
// It blocks until the provided context is canceled, then attempts to shut down
// the server within a 5-second grace period, and finally signals completion
// by closing the serverStopped channel.
func shutdownHandler(ctx context.Context, server *Server, serverStopped chan<- struct{}) {
	<-ctx.Done()
	server.logger.Info("shutdown signal received, initiating server shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		server.logger.Error("server shutdown error", zap.Error(err))
	} else {
		server.logger.Info("server HTTP shut down cleanly via Shutdown()")
	}

	close(serverStopped)
}

func run(ctx context.Context) error {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Sync()

	configPath, err := handleCLIWithFlagSet(flag.CommandLine)
	if err != nil {
		logger.Error("failed to parse CLI flags", zap.Error(err))
		return fmt.Errorf("failed to parse CLI flags: %w", err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	initMetrics()

	server := NewServer(cfg, logger)
	server.StartCollector(ctx)

	serverStopped := make(chan struct{})
	go shutdownHandler(ctx, server, serverStopped)

	sErr := server.Serve()

	<-serverStopped

	if sErr != nil && sErr != http.ErrServerClosed {
		logger.Error("server Serve error", zap.Error(sErr))
		return fmt.Errorf("server Serve error: %w", sErr)
	}

	logger.Info("Server goroutine finished.")
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	if err := run(ctx); err != nil {
		logger.Error("Application terminated with error", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Application exited cleanly.")
	os.Exit(0)
}
