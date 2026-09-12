package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"ttl-cli/i18n"
	clientapp "ttl-cli/internal/client/app"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"

	"github.com/spf13/cobra"
)

func TestWriteJSONSuccessUsesVersionedEnvelope(t *testing.T) {
	var output bytes.Buffer
	resource := clientapp.Resource{
		Key:   models.ValJsonKey{Key: "note", Type: models.ORIGIN},
		Value: models.ValJson{Val: "value", Tag: nil, CreatedAt: 1, UpdatedAt: 2},
	}

	if err := writeJSONSuccess(&output, resourceData{Resource: toResourceDTO(resource)}); err != nil {
		t.Fatalf("writeJSONSuccess() error = %v", err)
	}
	var envelope struct {
		SchemaVersion int  `json:"schema_version"`
		OK            bool `json:"ok"`
		Data          struct {
			Resource resourceDTO `json:"resource"`
		} `json:"data"`
	}
	if err := json.Unmarshal(output.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if envelope.SchemaVersion != 1 || !envelope.OK {
		t.Fatalf("envelope = %+v", envelope)
	}
	if envelope.Data.Resource.Tags == nil || len(envelope.Data.Resource.Tags) != 0 {
		t.Fatalf("tags = %#v, want []", envelope.Data.Resource.Tags)
	}
}

func TestWriteJSONErrorSortsAmbiguousCandidates(t *testing.T) {
	var output bytes.Buffer
	serviceErr := &clientapp.ServiceError{
		Kind:       clientapp.ErrorAmbiguous,
		Message:    "ambiguous",
		Candidates: []string{"z", "a"},
	}
	if err := writeJSONError(&output, serviceErr); err != nil {
		t.Fatalf("writeJSONError() error = %v", err)
	}
	var envelope errorEnvelope
	if err := json.Unmarshal(output.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	got, ok := envelope.Error.Details["candidates"].([]any)
	if !ok || !reflect.DeepEqual(got, []any{"a", "z"}) {
		t.Fatalf("candidates = %#v", envelope.Error.Details["candidates"])
	}
	if exitCodeFor(serviceErr, true) != exitConflict {
		t.Fatalf("exitCodeFor() = %d, want %d", exitCodeFor(serviceErr, true), exitConflict)
	}
}

func TestExitCodeForPreservesLegacyTextErrorCode(t *testing.T) {
	err := invalidArgument("bad input")
	if got := exitCodeFor(err, false); got != exitSystemError {
		t.Fatalf("exitCodeFor(text) = %d, want %d", got, exitSystemError)
	}
	if got := exitCodeFor(err, true); got != exitInvalidArgument {
		t.Fatalf("exitCodeFor(machine) = %d, want %d", got, exitInvalidArgument)
	}
}

func TestBoolFlagEnabledHandlesExplicitValues(t *testing.T) {
	for _, test := range []struct {
		args []string
		want bool
	}{
		{args: []string{"get", "--json"}, want: true},
		{args: []string{"get", "--json=true"}, want: true},
		{args: []string{"get", "--json=false"}, want: false},
		{args: []string{"add", "--", "--json", "value"}, want: false},
	} {
		if got := boolFlagEnabled(test.args, "--json"); got != test.want {
			t.Errorf("boolFlagEnabled(%v) = %v, want %v", test.args, got, test.want)
		}
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }

func TestReadValuePreservesStdinAndReportsFailure(t *testing.T) {
	root := newAddCommand(&options{})
	root.SetContext(context.Background())
	root.SetIn(bytes.NewBufferString("line 1\nline 2\n"))
	got, err := readValue(root, "-")
	if err != nil || got != "line 1\nline 2\n" {
		t.Fatalf("readValue() = %q, %v", got, err)
	}

	root.SetIn(failingReader{})
	if _, err := readValue(root, "-"); err == nil || mapCLIError(err).code != "system_error" {
		t.Fatalf("readValue() error = %v", err)
	}
}

func TestJSONSuccessIsDiscardedWhenCloseFails(t *testing.T) {
	var stdout, stderr bytes.Buffer
	opts := &options{service: clientapp.NewService(&closeErrorStorage{})}
	root := newRootCommand(opts)
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetIn(bytes.NewBuffer(nil))
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		opts.json = true
		cmd.SetContext(clientapp.WithService(context.WithValue(cmd.Context(), machineModeKey{}, invocationMode{json: true, nonInteractive: true}), opts.service))
		return nil
	}

	result := executeRoot(root, opts, []string{"get", "--json"}, true, &stdout)
	if len(result.stdout) != 0 || result.exitCode != exitSystemError {
		t.Fatalf("result = %+v", result)
	}
	var envelope errorEnvelope
	if err := json.Unmarshal(stderr.Bytes(), &envelope); err != nil || envelope.Error.Code != "system_error" {
		t.Fatalf("stderr = %q, error = %v", stderr.String(), err)
	}
}

func TestTextCommandErrorPreservesLocalizedCompatibility(t *testing.T) {
	if err := i18n.InitWithLanguage("en-US"); err != nil {
		t.Fatalf("i18n.InitWithLanguage() error = %v", err)
	}
	t.Cleanup(i18n.Reset)

	root := newAddCommand(&options{})
	root.SetContext(context.Background())
	err := &clientapp.ServiceError{Kind: clientapp.ErrorConflict, Message: "resource already exists: note"}
	if got := textCommandError(root, "add", "note", err).Error(); got != "Key already exists: note" {
		t.Fatalf("textCommandError(add) = %q", got)
	}

	err = &clientapp.ServiceError{Kind: clientapp.ErrorNotFound, Message: "resource not found: note"}
	mapped := textCommandError(root, "update", "note", err)
	if got := mapped.Error(); got != "Resource not found: note" {
		t.Fatalf("textCommandError(update) = %q", got)
	}
	if got := exitCodeFor(mapped, true); got != exitNotFound {
		t.Fatalf("exitCodeFor(non-interactive text) = %d, want %d", got, exitNotFound)
	}

	cause := errors.New("disk failed")
	err = &clientapp.ServiceError{Kind: clientapp.ErrorSystem, Operation: clientapp.ErrorRead, Message: "failed to read resources", Err: cause}
	mapped = textCommandError(root, "get", "", err)
	if got := mapped.Error(); got != "Failed to get resources: disk failed" || !errors.Is(mapped, cause) {
		t.Fatalf("textCommandError(get system) = %q, unwrap=%v", got, errors.Is(mapped, cause))
	}
}

type closeErrorStorage struct {
	corestorage.Storage
}

func (*closeErrorStorage) Close() error { return errors.New("close failed") }
func (*closeErrorStorage) GetAllResources() (map[models.ValJsonKey]models.ValJson, error) {
	return map[models.ValJsonKey]models.ValJson{}, nil
}
