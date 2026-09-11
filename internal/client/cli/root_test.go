package cli

import (
	"os"
	"reflect"
	"testing"
)

func TestNewRootCommand_HasClientAndCompatibilityCommands(t *testing.T) {
	root := NewRootCommand()
	for _, name := range []string{"add", "get", "sync", "migrate", "workspace", "server"} {
		cmd, _, err := root.Find([]string{name})
		if err != nil || cmd == nil || cmd.Name() != name {
			t.Fatalf("root command does not expose %s: cmd=%v err=%v", name, cmd, err)
		}
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

func TestServerArgs(t *testing.T) {
	original := os.Args
	t.Cleanup(func() { os.Args = original })
	os.Args = []string{"ttl", "--conf", "test.ini", "server", "user", "list"}
	got := serverArgs(newServerCompatibilityCommand())
	want := []string{"user", "list"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("serverArgs() = %v, want %v", got, want)
	}
}
