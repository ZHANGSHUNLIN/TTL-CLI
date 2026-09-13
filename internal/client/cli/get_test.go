package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	clientapp "ttl-cli/internal/client/app"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"

	"github.com/spf13/cobra"
)

type getStorage struct {
	corestorage.Storage
	resources map[models.ValJsonKey]models.ValJson
}

func (s *getStorage) Close() error { return nil }

func (s *getStorage) GetAllResources() (map[models.ValJsonKey]models.ValJson, error) {
	result := make(map[models.ValJsonKey]models.ValJson, len(s.resources))
	for key, value := range s.resources {
		result[key] = value
	}
	return result, nil
}

func (s *getStorage) SaveAuditRecord(models.AuditRecord) error { return nil }

func newGetTestCommand(storage *getStorage) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	service := clientapp.NewService(storage)
	cmd := newGetCommand(&options{service: service})
	ctx := context.WithValue(context.Background(), machineModeKey{}, invocationMode{})
	cmd.SetContext(clientapp.WithService(ctx, service))
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.SetIn(strings.NewReader(""))
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	return cmd, stdout, stderr
}

func TestGetCommand_DefaultDoesNotSearchValue(t *testing.T) {
	storage := &getStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "deployment-note", Type: models.ORIGIN}: {Val: "contains-secret"},
	}}
	cmd, stdout, _ := newGetTestCommand(storage)
	cmd.SetArgs([]string{"secret"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("get matched value without --value")
	}
	if strings.Contains(stdout.String(), "contains-secret") {
		t.Fatalf("stdout unexpectedly contains value: %q", stdout.String())
	}
}

func TestGetCommand_DefaultSearchesKeyAndTags(t *testing.T) {
	storage := &getStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "deployment-note", Type: models.ORIGIN}: {Val: "contains-secret", Tag: []string{"production"}},
	}}
	cmd, stdout, _ := newGetTestCommand(storage)
	cmd.SetArgs([]string{"production"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("get tag error = %v", err)
	}
	if !strings.Contains(stdout.String(), "contains-secret") {
		t.Fatalf("stdout = %q, want tag-matched resource value", stdout.String())
	}
}

func TestGetCommand_ValueFlagIncludesValueSearch(t *testing.T) {
	storage := &getStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "deployment-note", Type: models.ORIGIN}: {Val: "contains-secret"},
	}}
	cmd, stdout, _ := newGetTestCommand(storage)
	cmd.SetArgs([]string{"--value", "secret"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("get --value error = %v", err)
	}
	if !strings.Contains(stdout.String(), "contains-secret") {
		t.Fatalf("stdout = %q, want matched value", stdout.String())
	}
	if cmd.Flags().Lookup("value") == nil {
		t.Fatal("get command does not expose --value")
	}
}

func TestGetCommand_ValueFlagSupportsShortAliases(t *testing.T) {
	storage := &getStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "deployment-note", Type: models.ORIGIN}: {Val: "contains-secret"},
	}}
	cmd, stdout, _ := newGetTestCommand(storage)
	cmd.SetArgs([]string{"-v", "secret"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("get -v error = %v", err)
	}
	if !strings.Contains(stdout.String(), "contains-secret") {
		t.Fatalf("stdout = %q, want matched value", stdout.String())
	}
	if cmd.Flags().ShorthandLookup("v") == nil {
		t.Fatal("get command does not expose -v")
	}
}

func TestExecuteRoot_NormalizesGetValueAlias(t *testing.T) {
	storage := &getStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "deployment-note", Type: models.ORIGIN}: {Val: "contains-secret"},
	}}
	service := clientapp.NewService(storage)
	opts := &options{service: service}
	root := newRootCommand(opts)
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		ctx := context.WithValue(cmd.Context(), machineModeKey{}, invocationMode{})
		cmd.SetContext(clientapp.WithService(ctx, service))
		return nil
	}
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetIn(strings.NewReader(""))

	result := executeRoot(root, opts, []string{"get", "-val", "secret"}, false, &bytes.Buffer{})
	if result.exitCode != exitSuccess {
		t.Fatalf("get -val exit code = %d, stderr = %q", result.exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "contains-secret") {
		t.Fatalf("stdout = %q, want matched value", stdout.String())
	}
}
