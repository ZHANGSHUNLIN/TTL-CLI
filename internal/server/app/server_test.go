package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateConfig(t *testing.T) {
	base := Config{ListenAddress: "127.0.0.1", Port: 8080, DataDir: t.TempDir(), ShutdownTimeout: time.Second}
	tests := []struct {
		name string
		edit func(*Config)
	}{
		{name: "empty listen address", edit: func(cfg *Config) { cfg.ListenAddress = "" }},
		{name: "invalid low port", edit: func(cfg *Config) { cfg.Port = 0 }},
		{name: "invalid high port", edit: func(cfg *Config) { cfg.Port = 65536 }},
		{name: "empty data directory", edit: func(cfg *Config) { cfg.DataDir = "" }},
		{name: "non-positive shutdown timeout", edit: func(cfg *Config) { cfg.ShutdownTimeout = 0 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := base
			tt.edit(&cfg)
			if err := validateConfig(cfg); err == nil {
				t.Fatal("validateConfig returned nil")
			}
		})
	}
}

func TestValidateDataDirWritable(t *testing.T) {
	if err := validateDataDirWritable(t.TempDir()); err != nil {
		t.Fatalf("validateDataDirWritable returned error: %v", err)
	}
}

func TestRun_StartsHealthEndpointAndStopsOnContext(t *testing.T) {
	port := freePort(t)
	cfg := Config{
		ListenAddress:   "127.0.0.1",
		Port:            port,
		DataDir:         t.TempDir(),
		ShutdownTimeout: time.Second,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- Run(ctx, cfg) }()

	url := fmt.Sprintf("http://127.0.0.1:%d/healthz", port)
	var response *http.Response
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		candidate, err := http.Get(url)
		if err == nil {
			response = candidate
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if response == nil {
		t.Fatal("health endpoint did not become reachable")
	}
	if response.Body != nil {
		_ = response.Body.Close()
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	cancel()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestRun_FailsBeforeListeningWhenUsersFileIsCorrupt(t *testing.T) {
	dataDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dataDir, "users.json"), []byte("{"), 0600); err != nil {
		t.Fatalf("write corrupt users file: %v", err)
	}

	err := Run(context.Background(), Config{
		ListenAddress:   "127.0.0.1",
		Port:            freePort(t),
		DataDir:         dataDir,
		ShutdownTimeout: time.Second,
	})
	if err == nil {
		t.Fatal("Run returned nil for corrupt users.json")
	}
}

func TestRun_FailsWhenDataDirPathIsNotDirectory(t *testing.T) {
	dataPath := filepath.Join(t.TempDir(), "data")
	if err := os.WriteFile(dataPath, []byte("not a directory"), 0600); err != nil {
		t.Fatalf("write data path: %v", err)
	}

	err := Run(context.Background(), Config{
		ListenAddress:   "127.0.0.1",
		Port:            freePort(t),
		DataDir:         dataPath,
		ShutdownTimeout: time.Second,
	})
	if err == nil {
		t.Fatal("Run returned nil when data directory path was a file")
	}
}

func TestRun_FailsWhenPortIsOccupied(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("occupy port: %v", err)
	}
	defer occupied.Close()

	port := occupied.Addr().(*net.TCPAddr).Port
	err = Run(context.Background(), Config{
		ListenAddress:   "127.0.0.1",
		Port:            port,
		DataDir:         t.TempDir(),
		ShutdownTimeout: time.Second,
	})
	if err == nil {
		t.Fatal("Run returned nil for occupied port")
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("allocate port: %v", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
