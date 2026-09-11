package cli

import (
	"os"

	"ttl-cli/i18n"
	servercli "ttl-cli/internal/server/cli"

	"github.com/spf13/cobra"
)

func newServerCompatibilityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "server",
		Short:              i18n.T("command.server.short"),
		Long:               i18n.T("command.server.long"),
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			server := servercli.NewRootCommand()
			server.SetArgs(serverArgs(command))
			server.SetIn(os.Stdin)
			server.SetOut(os.Stdout)
			server.SetErr(os.Stderr)
			return server.Execute()
		},
	}
	cmd.SetHelpFunc(func(command *cobra.Command, _ []string) {
		_, _ = command.OutOrStdout().Write([]byte("Run 'ttl-server --help' for server commands.\n"))
	})
	return cmd
}

func serverArgs(command *cobra.Command) []string {
	for i, arg := range os.Args {
		if arg == command.Name() {
			return os.Args[i+1:]
		}
	}
	return nil
}
