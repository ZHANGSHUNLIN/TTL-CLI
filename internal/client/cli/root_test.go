package cli

import (
	"testing"
)

func TestNewRootCommand_HasClientCommands(t *testing.T) {
	root := NewRootCommand()
	for _, name := range []string{"add", "get", "sync", "migrate", "workspace"} {
		cmd, _, err := root.Find([]string{name})
		if err != nil || cmd == nil || cmd.Name() != name {
			t.Fatalf("root command does not expose %s: cmd=%v err=%v", name, cmd, err)
		}
	}
	if _, _, err := root.Find([]string{"server"}); err == nil {
		t.Fatal("client root must not expose server compatibility command")
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
