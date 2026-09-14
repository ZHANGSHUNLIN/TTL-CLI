package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	clientapp "ttl-cli/internal/client/app"
	clienttui "ttl-cli/internal/client/tui"
	corestorage "ttl-cli/internal/core/storage"
)

func TestNewRootCommand_HasClientCommands(t *testing.T) {
	root := NewRootCommand()
	for _, name := range []string{"add", "get", "pick", "sync", "workspace", "ui"} {
		cmd, _, err := root.Find([]string{name})
		if err != nil || cmd == nil || cmd.Name() != name {
			t.Fatalf("root command does not expose %s: cmd=%v err=%v", name, cmd, err)
		}
	}
	if _, _, err := root.Find([]string{"migrate"}); err == nil {
		t.Fatal("client root must not expose removed migrate command")
	}
	if _, _, err := root.Find([]string{"server"}); err == nil {
		t.Fatal("client root must not expose server compatibility command")
	}
}

func TestUICommand_GatesUnsupportedStorageAndNonTerminalBeforeOpen(t *testing.T) {
	for _, test := range []struct {
		name      string
		storage   string
		terminal  bool
		wantError string
	}{
		{name: "cloud", storage: "cloud", terminal: true, wantError: "仅支持 local"},
		{name: "sync", storage: "sync", terminal: true, wantError: "不支持的存储模式"},
		{name: "non terminal", storage: "local", terminal: false, wantError: "interactive terminal"},
	} {
		t.Run(test.name, func(t *testing.T) {
			openCalls := 0
			opts := &options{storageType: test.storage}
			opts.isTerminal = func(io.Reader, io.Writer) bool { return test.terminal }
			opts.openStorage = func(string, string, string, int, string) (corestorage.Storage, error) {
				openCalls++
				return nil, errors.New("must not open")
			}
			root := newRootCommand(opts)
			root.SetIn(bytes.NewBuffer(nil))
			root.SetOut(&bytes.Buffer{})
			root.SetErr(&bytes.Buffer{})
			root.SetArgs([]string{"--storage", test.storage, "ui"})
			err := root.Execute()
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("Execute() error = %v, want %q", err, test.wantError)
			}
			if openCalls != 0 {
				t.Fatalf("open calls = %d, want 0", openCalls)
			}
		})
	}
}

func TestUICommand_RunsWithInjectedService(t *testing.T) {
	storage := &closeErrorStorage{}
	service := clientapp.NewService(storage)
	runCalls := 0
	opts := &options{service: service}
	opts.runTUI = func(got clienttui.ResourceService, _ clienttui.RunOptions) error {
		runCalls++
		if got != service {
			t.Fatalf("runner service = %T, want injected service", got)
		}
		return nil
	}
	cmd := newUICommand(opts)
	cmd.SetContext(clientapp.WithService(context.Background(), service))
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if runCalls != 1 {
		t.Fatalf("runner calls = %d, want 1", runCalls)
	}
}

func TestCountSpecialChars(t *testing.T) {
	t.Parallel()
	tests := map[string]int{"": 0, "~": 1, "^^": 2, "~^": 2, "~x": 0}
	for input, want := range tests {
		if got := countSpecialChars(input); got != want {
			t.Errorf("countSpecialChars(%q) = %d, want %d", input, got, want)
		}
	}
}
