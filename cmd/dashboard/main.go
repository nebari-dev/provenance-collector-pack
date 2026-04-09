package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nebari-dev/provenance-collector/internal/dashboard"
	"github.com/nebari-dev/provenance-collector/internal/report"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	slog.Info("starting provenance dashboard", "version", report.Version)

	reportsDir := os.Getenv("PROVENANCE_REPORT_PATH")
	if reportsDir == "" {
		reportsDir = "/reports"
	}

	addr := os.Getenv("PROVENANCE_DASHBOARD_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	slog.Info("configuration loaded",
		"addr", addr,
		"reportsDir", reportsDir,
	)

	srv := dashboard.NewServer(reportsDir)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("starting dashboard", "addr", addr, "reports", reportsDir)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down dashboard")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
