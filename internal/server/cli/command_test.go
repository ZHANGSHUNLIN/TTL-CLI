package cli

import (
	"testing"

	"ttl-cli/internal/server/tenant"

	"github.com/spf13/cobra"
)

func TestNewRootCommand_HasSeparateServeAndUserCommands(t *testing.T) {
	t.Parallel()

	root := NewRootCommand()
	if root.Use != "ttl-server" {
		t.Fatalf("root Use = %q, want ttl-server", root.Use)
	}
	for _, name := range []string{"listen", "port", "data-dir", "shutdown-timeout"} {
		if root.PersistentFlags().Lookup(name) == nil {
			t.Fatalf("standalone server command must expose --%s", name)
		}
	}
	if commandNamed(root, "serve") == nil {
		t.Fatal("standalone server command must expose serve")
	}
	user := commandNamed(root, "user")
	if user == nil {
		t.Fatal("standalone server command must expose user")
	}
	for _, name := range []string{"add", "list", "enable", "disable", "reset-key", "delete"} {
		if commandNamed(user, name) == nil {
			t.Fatalf("user command must expose %s", name)
		}
	}
}

func TestNewCompatibilityCommand_PreservesLegacyServerShape(t *testing.T) {
	t.Parallel()

	server := NewCompatibilityCommand()
	if server.Use != "server" {
		t.Fatalf("compatibility Use = %q, want server", server.Use)
	}
	if server.Flags().Lookup("port") == nil || server.Flags().Lookup("data-dir") == nil {
		t.Fatal("legacy server command must preserve port and data-dir flags")
	}
	if commandNamed(server, "user") == nil {
		t.Fatal("legacy server command must preserve user subcommand")
	}
}

func TestNewRootCommand_UserLifecycle(t *testing.T) {
	dataDir := t.TempDir()

	executeCommand(t, NewRootCommand(),
		"--data-dir", dataDir,
		"user", "add",
		"--id", "alice",
		"--name", "Alice",
	)
	executeCommand(t, NewRootCommand(), "--data-dir", dataDir, "user", "list")

	executeCommand(t, NewRootCommand(), "--data-dir", dataDir, "user", "disable", "--id", "alice")
	executeCommand(t, NewRootCommand(), "--data-dir", dataDir, "user", "enable", "--id", "alice")
	executeCommand(t, NewRootCommand(), "--data-dir", dataDir, "user", "reset-key", "--id", "alice")
	executeCommand(t, NewRootCommand(), "--data-dir", dataDir, "user", "delete", "--id", "alice", "--confirm")

	store := tenant.NewUserStore(dataDir + "/users.json")
	if err := store.Load(); err != nil {
		t.Fatalf("load user store: %v", err)
	}
	if users := store.ListUsers(); len(users) != 0 {
		t.Fatalf("users after delete = %d, want 0", len(users))
	}
}

func executeCommand(t *testing.T, root *cobra.Command, args ...string) {
	t.Helper()

	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		t.Fatalf("execute %v: %v", args, err)
	}
}

func commandNamed(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
