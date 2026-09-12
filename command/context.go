package command

import (
	"github.com/spf13/cobra"
	clientapp "ttl-cli/internal/client/app"
)

// clientService returns the service scoped to the current command execution.
// A nil service is represented by an empty Service so handlers return the
// existing "storage not initialized" error instead of panicking in tests.
func clientService(cmd *cobra.Command) *clientapp.Service {
	service, err := clientapp.ServiceFromContext(cmd.Context())
	if err != nil {
		return clientapp.NewService(nil)
	}
	return service
}
