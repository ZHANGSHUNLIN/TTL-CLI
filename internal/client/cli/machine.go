package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	clientapp "ttl-cli/internal/client/app"
	"ttl-cli/models"

	"github.com/spf13/cobra"
)

const (
	exitSuccess         = 0
	exitSystemError     = 1
	exitInvalidArgument = 2
	exitNotFound        = 3
	exitConflict        = 4
)

type invocationMode struct {
	json           bool
	nonInteractive bool
}

type machineModeKey struct{}

type successEnvelope struct {
	SchemaVersion int  `json:"schema_version"`
	OK            bool `json:"ok"`
	Data          any  `json:"data"`
}

type errorEnvelope struct {
	SchemaVersion int          `json:"schema_version"`
	OK            bool         `json:"ok"`
	Error         machineError `json:"error"`
}

type machineError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type resourceDTO struct {
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Tags      []string `json:"tags"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
}

type resourceData struct {
	Resource resourceDTO `json:"resource"`
}

type resourcesData struct {
	Resources []resourceDTO `json:"resources"`
}

type deleteData struct {
	Key     string `json:"key"`
	Deleted bool   `json:"deleted"`
}

type cliError struct {
	code     string
	message  string
	details  map[string]any
	exitCode int
	cause    error
}

func (e *cliError) Error() string { return e.message }
func (e *cliError) Unwrap() error { return e.cause }

func modeFromCommand(cmd *cobra.Command) invocationMode {
	if mode, ok := cmd.Context().Value(machineModeKey{}).(invocationMode); ok {
		return mode
	}
	jsonMode, _ := cmd.Flags().GetBool("json")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	return invocationMode{json: jsonMode, nonInteractive: nonInteractive || jsonMode}
}

func validateMachineMode(cmd *cobra.Command) error {
	mode := modeFromCommand(cmd)
	if !mode.json && !mode.nonInteractive {
		return nil
	}
	if machineCommandSupported(cmd.Name()) {
		return nil
	}
	return invalidArgument(fmt.Sprintf("machine mode is not supported for command %q", cmd.CommandPath()))
}

func machineCommandSupported(name string) bool {
	switch name {
	case "add", "get", "update", "put", "mod", "del", "rm", "tag", "dtag":
		return true
	default:
		return false
	}
}

func writeJSONSuccess(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(successEnvelope{SchemaVersion: 1, OK: true, Data: data})
}

func writeJSONError(w io.Writer, err error) error {
	mapped := mapCLIError(err)
	details := mapped.details
	if details == nil {
		details = map[string]any{}
	}
	return json.NewEncoder(w).Encode(errorEnvelope{
		SchemaVersion: 1,
		OK:            false,
		Error: machineError{
			Code:    mapped.code,
			Message: mapped.message,
			Details: details,
		},
	})
}

func exitCodeFor(err error, machineMode bool) int {
	if err == nil {
		return exitSuccess
	}
	if !machineMode {
		return exitSystemError
	}
	return mapCLIError(err).exitCode
}

func mapCLIError(err error) *cliError {
	var known *cliError
	if errors.As(err, &known) {
		return known
	}

	var serviceErr *clientapp.ServiceError
	if errors.As(err, &serviceErr) {
		details := map[string]any{}
		if len(serviceErr.Candidates) > 0 {
			candidates := append([]string(nil), serviceErr.Candidates...)
			sort.Strings(candidates)
			details["candidates"] = candidates
		}
		switch serviceErr.Kind {
		case clientapp.ErrorNotFound:
			return &cliError{code: "not_found", message: serviceErr.Error(), details: details, exitCode: exitNotFound, cause: err}
		case clientapp.ErrorConflict:
			return &cliError{code: "conflict", message: serviceErr.Error(), details: details, exitCode: exitConflict, cause: err}
		case clientapp.ErrorAmbiguous:
			return &cliError{code: "ambiguous", message: serviceErr.Error(), details: details, exitCode: exitConflict, cause: err}
		}
	}
	return &cliError{code: "system_error", message: err.Error(), details: map[string]any{}, exitCode: exitSystemError, cause: err}
}

func invalidArgument(message string) error {
	return &cliError{code: "invalid_argument", message: message, details: map[string]any{}, exitCode: exitInvalidArgument}
}

func interactionRequired(message string) error {
	return &cliError{code: "interaction_required", message: message, details: map[string]any{}, exitCode: exitConflict}
}

func isArgumentError(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.HasPrefix(message, "unknown flag:") ||
		strings.HasPrefix(message, "unknown command") ||
		strings.Contains(message, "requires at least") ||
		strings.Contains(message, "accepts ") ||
		strings.Contains(message, "requires exactly") ||
		strings.Contains(message, "invalid argument")
}

func ambiguousError(matches []clientapp.Resource) error {
	candidates := make([]string, 0, len(matches))
	for _, match := range matches {
		candidates = append(candidates, resourceKey(match.Key))
	}
	return &clientapp.ServiceError{
		Kind:       clientapp.ErrorAmbiguous,
		Message:    "multiple resources matched the query",
		Candidates: candidates,
	}
}

func readValue(cmd *cobra.Command, value string) (string, error) {
	if value != "-" {
		return value, nil
	}
	content, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return "", &cliError{code: "system_error", message: "failed to read value from stdin", details: map[string]any{}, exitCode: exitSystemError, cause: err}
	}
	return string(content), nil
}

func toResourceDTO(resource clientapp.Resource) resourceDTO {
	tags := append([]string(nil), resource.Value.Tag...)
	if tags == nil {
		tags = []string{}
	}
	return resourceDTO{
		Key:       resourceKey(resource.Key),
		Value:     resource.Value.Val,
		Tags:      tags,
		CreatedAt: resource.Value.CreatedAt,
		UpdatedAt: resource.Value.UpdatedAt,
	}
}

func resourceKey(key models.ValJsonKey) string {
	if key.Type == models.TAG && key.OriginKey != "" {
		return key.OriginKey
	}
	return key.Key
}
