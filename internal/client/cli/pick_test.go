package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	clientapp "ttl-cli/internal/client/app"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"

	"github.com/spf13/cobra"
)

type pickStorage struct {
	corestorage.Storage
	resources map[models.ValJsonKey]models.ValJson
}

func (s *pickStorage) GetAllResources() (map[models.ValJsonKey]models.ValJson, error) {
	result := make(map[models.ValJsonKey]models.ValJson, len(s.resources))
	for key, value := range s.resources {
		result[key] = value
	}
	return result, nil
}

func newPickTestCommand(storage *pickStorage, input io.Reader, terminal bool) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	service := clientapp.NewService(storage)
	opts := &options{
		service: service,
		isTerminal: func(io.Reader, io.Writer) bool {
			return terminal
		},
	}
	cmd := newPickCommand(opts)
	cmd.SilenceUsage = true
	ctx := context.WithValue(context.Background(), machineModeKey{}, invocationMode{})
	cmd.SetContext(clientapp.WithService(ctx, service))
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.SetIn(input)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	return cmd, stdout, stderr
}

func TestPickCommand_SingleMatchKeepsStdoutCleanAndNormalizesNewline(t *testing.T) {
	storage := &pickStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "note", Type: models.ORIGIN}: {Val: "value\n\n", CreatedAt: 1},
	}}
	cmd, stdout, stderr := newPickTestCommand(storage, &readFailReader{}, false)
	cmd.SetArgs([]string{"note"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := stdout.String(); got != "value\n" {
		t.Fatalf("stdout = %q, want one terminating newline", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestPickCommand_MultipleMatchesWritesCandidatesToStderr(t *testing.T) {
	storage := &pickStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "older", Type: models.ORIGIN}: {Val: "old", CreatedAt: 1},
		{Key: "newer", Type: models.ORIGIN}: {Val: "new", CreatedAt: 2},
	}}
	cmd, stdout, stderr := newPickTestCommand(storage, strings.NewReader("2\n"), true)
	cmd.SetArgs([]string{""})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stdout.String() != "old\n" {
		t.Fatalf("stdout = %q, want selected value only", stdout.String())
	}
	if got := stderr.String(); !strings.Contains(got, "1. newer") || !strings.Contains(got, "2. older") {
		t.Fatalf("stderr = %q, want stable candidate list", got)
	}
	if strings.Contains(stdout.String(), "newer") || strings.Contains(stdout.String(), "older") {
		t.Fatalf("stdout contains candidate diagnostics: %q", stdout.String())
	}
}

func TestPickCommand_InvalidChoiceAndCancellationHaveStableCodes(t *testing.T) {
	tests := []struct {
		name string
		in   string
		code string
	}{
		{name: "invalid", in: "9\n", code: "invalid_choice"},
		{name: "cancel", in: "q\n", code: "cancelled"},
		{name: "eof", in: "", code: "cancelled"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := &pickStorage{resources: map[models.ValJsonKey]models.ValJson{
				{Key: "a", Type: models.ORIGIN}: {Val: "a", CreatedAt: 1},
				{Key: "b", Type: models.ORIGIN}: {Val: "b", CreatedAt: 2},
			}}
			cmd, stdout, _ := newPickTestCommand(storage, strings.NewReader(test.in), true)
			cmd.SetArgs([]string{})
			err := cmd.Execute()
			if err == nil || mapCLIError(err).code != test.code {
				t.Fatalf("error = %v, code = %q, want %q", err, mapCLIError(err).code, test.code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
		})
	}
}

func TestPickCommand_NoTTYFailsBeforeReadingInput(t *testing.T) {
	storage := &pickStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "a", Type: models.ORIGIN}: {Val: "a", CreatedAt: 1},
		{Key: "b", Type: models.ORIGIN}: {Val: "b", CreatedAt: 2},
	}}
	cmd, stdout, _ := newPickTestCommand(storage, &readFailReader{}, false)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil || mapCLIError(err).code != "interaction_required" {
		t.Fatalf("error = %v, code = %q, want interaction_required", err, mapCLIError(err).code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func TestPickCommand_EmptyResultIsNotFound(t *testing.T) {
	cmd, stdout, _ := newPickTestCommand(&pickStorage{resources: map[models.ValJsonKey]models.ValJson{}}, strings.NewReader(""), false)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil || mapCLIError(err).code != "not_found" {
		t.Fatalf("error = %v, code = %q, want not_found", err, mapCLIError(err).code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func TestShouldRecordHistory_PickIsReadOnly(t *testing.T) {
	cmd := &cobra.Command{Use: "pick"}
	if shouldRecordHistory(cmd) {
		t.Fatal("pick should not record command history")
	}
}

func TestPickCommand_RejectsMachineModes(t *testing.T) {
	for _, flag := range []string{"--json", "--non-interactive"} {
		t.Run(flag, func(t *testing.T) {
			openCalls := 0
			opts := &options{openStorage: func(string, string, string, int, string) (corestorage.Storage, error) {
				openCalls++
				return nil, errors.New("storage must not open")
			}}
			root := newRootCommand(opts)
			root.SetIn(strings.NewReader(""))
			root.SetOut(&bytes.Buffer{})
			root.SetErr(&bytes.Buffer{})
			root.SetArgs([]string{flag, "pick"})
			err := root.Execute()
			if err == nil || mapCLIError(err).code != "invalid_argument" {
				t.Fatalf("Execute() error = %v, code = %q, want invalid_argument", err, mapCLIError(err).code)
			}
			if openCalls != 0 {
				t.Fatalf("open calls = %d, want 0", openCalls)
			}
		})
	}
}

type readFailReader struct{}

func (*readFailReader) Read([]byte) (int, error) { return 0, errors.New("input must not be read") }
