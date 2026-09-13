package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"ttl-cli/internal/server/api"
	"ttl-cli/internal/server/tenant"
)

// Config defines the standalone server process contract.
type Config struct {
	ListenAddress   string
	Port            int
	DataDir         string
	ShutdownTimeout time.Duration
}

// RunWithSignals runs the server until SIGINT, SIGTERM, or a server error.
func RunWithSignals(cfg Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	return Run(ctx, cfg)
}

// Run starts the HTTP server and closes all tenant storage on exit.
func Run(ctx context.Context, cfg Config) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.DataDir, 0700); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	if err := validateDataDirWritable(cfg.DataDir); err != nil {
		return err
	}

	userStore := tenant.NewUserStore(filepath.Join(cfg.DataDir, "users.json"))
	if err := userStore.Load(); err != nil {
		return fmt.Errorf("failed to load user data: %w", err)
	}
	users := userStore.ListUsers()
	if len(users) == 0 {
		fmt.Println("Warning: No users, please run 'ttl-server user add' to create users first")
	}

	tenantMgr := tenant.NewStorageManager(filepath.Join(cfg.DataDir, "tenants"))
	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.ListenAddress, fmt.Sprintf("%d", cfg.Port)))
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%d: %w", cfg.ListenAddress, cfg.Port, err)
	}

	var ready atomic.Bool
	server := &http.Server{
		Handler:           api.NewServerHandler(userStore, tenantMgr, ready.Load),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ready.Store(true)
	fmt.Printf("ttl-server started, listening on %s\n", listener.Addr())
	fmt.Printf("  REST API: http://%s/api/v1/\n", listener.Addr())
	fmt.Printf("  Health: http://%s/healthz\n", listener.Addr())
	fmt.Printf("  Data dir: %s\n", cfg.DataDir)
	fmt.Printf("  User count: %d\n", len(users))

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(listener)
	}()

	select {
	case err := <-serveErr:
		_ = tenantMgr.CloseAll()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		ready.Store(false)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		shutdownErr := server.Shutdown(shutdownCtx)
		cancel()
		closeErr := tenantMgr.CloseAll()
		if shutdownErr != nil {
			return fmt.Errorf("server shutdown failed: %w", shutdownErr)
		}
		if closeErr != nil {
			return closeErr
		}
		return nil
	}
}

func validateConfig(cfg Config) error {
	if cfg.ListenAddress == "" {
		return errors.New("listen address cannot be empty")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535: %d", cfg.Port)
	}
	if cfg.DataDir == "" {
		return errors.New("data directory cannot be empty")
	}
	if cfg.ShutdownTimeout <= 0 {
		return errors.New("shutdown timeout must be positive")
	}
	return nil
}

func validateDataDirWritable(dataDir string) error {
	testFile, err := os.CreateTemp(dataDir, ".ttl-server-write-check-")
	if err != nil {
		return fmt.Errorf("data directory is not writable: %w", err)
	}
	name := testFile.Name()
	if err := testFile.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("data directory write check failed: %w", err)
	}
	if err := os.Remove(name); err != nil {
		return fmt.Errorf("data directory cleanup failed: %w", err)
	}
	return nil
}
