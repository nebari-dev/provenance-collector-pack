package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

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

	authCfg := dashboard.AuthConfig{
		IssuerURL:   os.Getenv("PROVENANCE_OIDC_ISSUER"),
		AdminGroups: splitAndTrim(os.Getenv("PROVENANCE_ADMIN_GROUPS")),
	}

	manualJobTTL := parseManualJobTTL(os.Getenv("PROVENANCE_MANUAL_JOB_TTL"))

	slog.Info("configuration loaded",
		"addr", addr,
		"reportsDir", reportsDir,
		"authIssuer", authCfg.IssuerURL,
		"adminGroups", len(authCfg.AdminGroups),
		"manualJobTTL", manualJobTTL.String(),
	)

	srv := dashboard.NewServer(reportsDir).WithAuth(authCfg)

	// /api/scan needs an in-cluster client + the CronJob's namespace/name.
	// Missing config or being out-of-cluster simply leaves the endpoint
	// disabled — handler will respond 503.
	namespace := os.Getenv("PROVENANCE_NAMESPACE")
	cronJobName := os.Getenv("PROVENANCE_CRONJOB_NAME")
	if namespace != "" && cronJobName != "" {
		if runner, err := buildScanRunner(namespace, cronJobName, manualJobTTL); err != nil {
			slog.Warn("scan endpoint disabled: kubernetes client unavailable",
				"namespace", namespace, "cronJob", cronJobName, "error", err)
		} else {
			srv = srv.WithScanRunner(runner)
			slog.Info("scan endpoint enabled",
				"namespace", namespace, "cronJob", cronJobName)
		}
	}

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

func buildScanRunner(namespace, cronJobName string, manualJobTTL time.Duration) (dashboard.ScanRunner, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return dashboard.NewK8sScanRunner(client, namespace, cronJobName, manualJobTTL), nil
}

// parseManualJobTTL reads the PROVENANCE_MANUAL_JOB_TTL env var.
// Empty / unparseable values fall back to DefaultManualJobTTL.
// Set the value to "0" (or any zero-duration string) to disable the TTL
// entirely so manual Jobs persist until manually deleted.
func parseManualJobTTL(v string) time.Duration {
	if v == "" {
		return dashboard.DefaultManualJobTTL
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		slog.Warn("invalid PROVENANCE_MANUAL_JOB_TTL, using default", "value", v, "error", err)
		return dashboard.DefaultManualJobTTL
	}
	if d < 0 {
		return 0
	}
	return d
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
