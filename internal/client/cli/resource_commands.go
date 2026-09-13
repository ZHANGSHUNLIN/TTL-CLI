package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"

	"ttl-cli/i18n"
	clientapp "ttl-cli/internal/client/app"
	"ttl-cli/models"
	"ttl-cli/util"

	"github.com/spf13/cobra"
)

func newAddCommand(_ *options) *cobra.Command {
	var tags []string
	cmd := &cobra.Command{
		Use:   "add [key] [value]",
		Short: i18n.T("command.add.short"),
		Long:  i18n.T("command.add.long"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := modeFromCommand(cmd)
			value, err := readValue(cmd, args[1])
			if err != nil {
				return err
			}
			if args[1] != "-" {
				value = util.UnescapeString(value)
			}
			resource, err := serviceFromCommand(cmd).CreateResource(args[0], value, tags)
			if err != nil {
				return textCommandError(cmd, "add", args[0], err)
			}
			debug, _ := cmd.Context().Value("debug").(bool)
			if err := serviceFromCommand(cmd).RecordAudit(args[0], "add"); err != nil && debug && !mode.json {
				fmt.Fprintf(cmd.OutOrStdout(), i18n.T("command.add.audit_error"), err)
			}
			if mode.json {
				return writeJSONSuccess(cmd.OutOrStdout(), resourceData{Resource: toResourceDTO(resource)})
			}
			fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.add.success"))
			return nil
		},
	}
	cmd.Flags().StringSliceVarP(&tags, "tag", "t", nil, i18n.T("command.add.flag_tag"))
	return cmd
}

func newGetCommand(_ *options) *cobra.Command {
	var includeValue bool
	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: i18n.T("command.get.short"),
		Long:  i18n.T("command.get.long"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := modeFromCommand(cmd)
			service := serviceFromCommand(cmd)
			if len(args) == 0 {
				resources, err := service.ListResources()
				if err != nil {
					return textCommandError(cmd, "get", "", err)
				}
				if mode.json {
					dtos := make([]resourceDTO, 0, len(resources))
					for _, resource := range resources {
						dtos = append(dtos, toResourceDTO(resource))
					}
					return writeJSONSuccess(cmd.OutOrStdout(), resourcesData{Resources: dtos})
				}
				fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.get.no_filter_notice"))
				fmt.Fprintln(cmd.OutOrStdout())
				for _, resource := range resources {
					fmt.Fprintln(cmd.OutOrStdout(), "  ", resource.Key.Key)
				}
				fmt.Fprintln(cmd.OutOrStdout())
				return nil
			}

			matches, err := service.FindResourcesWithOptions(args[0], clientapp.SearchOptions{IncludeValue: includeValue, IncludeTags: true})
			if err != nil {
				return textCommandError(cmd, "get", args[0], err)
			}
			var selected clientapp.Resource
			switch {
			case len(matches) == 1:
				selected = matches[0]
			case mode.nonInteractive:
				return ambiguousError(matches)
			default:
				selected, err = selectResource(cmd, matches)
				if err != nil {
					return err
				}
			}

			debug, _ := cmd.Context().Value("debug").(bool)
			if err := service.RecordAudit(resourceKey(selected.Key), "get"); err != nil && debug && !mode.json {
				fmt.Fprintf(cmd.OutOrStdout(), i18n.T("command.get.audit_error"), err)
			}
			if mode.json {
				return writeJSONSuccess(cmd.OutOrStdout(), resourceData{Resource: toResourceDTO(selected)})
			}
			printTextResource(cmd, selected)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&includeValue, "value", "v", false, i18n.T("command.get.flag_value"))
	// Keep the requested multi-character spelling as a hidden compatibility alias;
	// -v is the discoverable conventional shorthand.
	cmd.Flags().BoolVar(&includeValue, "val", false, i18n.T("command.get.flag_value"))
	_ = cmd.Flags().MarkHidden("val")
	return cmd
}

func newPickCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "pick [query]",
		Short: i18n.T("command.pick.short"),
		Long:  i18n.T("command.pick.long"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service := serviceFromCommand(cmd)
			var (
				matches []clientapp.Resource
				err     error
			)
			if len(args) == 0 {
				matches, err = service.ListResources()
			} else {
				matches, err = service.FindResources(args[0])
			}
			if err != nil {
				return err
			}
			if len(matches) == 0 {
				query := ""
				if len(args) == 1 {
					query = args[0]
				}
				return pickNotFoundError(query)
			}

			selected := matches[0]
			if len(matches) > 1 {
				detector := opts.isTerminal
				if detector == nil {
					detector = defaultTerminalDetector
				}
				if !detector(cmd.InOrStdin(), cmd.ErrOrStderr()) {
					return interactionRequired(i18n.T("command.pick.requires_terminal"))
				}
				var selectErr error
				selected, selectErr = selectPickResource(cmd, matches)
				if selectErr != nil {
					return selectErr
				}
			}
			writePickValue(cmd.OutOrStdout(), selected.Value.Val)
			return nil
		},
	}
}

func selectPickResource(cmd *cobra.Command, matches []clientapp.Resource) (clientapp.Resource, error) {
	fmt.Fprintln(cmd.ErrOrStderr(), i18n.T("command.pick.multiple_matches"))
	for index, match := range matches {
		fmt.Fprintf(cmd.ErrOrStderr(), "%d. %s\n", index+1, resourceKey(match.Key))
	}
	fmt.Fprintln(cmd.ErrOrStderr(), i18n.T("command.pick.prompt"))

	type readResult struct {
		input string
		err   error
	}
	readCh := make(chan readResult, 1)
	go func() {
		input, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
		readCh <- readResult{input: input, err: err}
	}()
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	defer signal.Stop(interruptCh)
	var input string
	var err error
	select {
	case result := <-readCh:
		input, err = result.input, result.err
	case <-interruptCh:
		return clientapp.Resource{}, pickCancelledError()
	}
	if len(input) > 0 {
		switch input[0] {
		case '\x03', '\x1b':
			return clientapp.Resource{}, pickCancelledError()
		}
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return clientapp.Resource{}, &cliError{code: "system_error", message: err.Error(), details: map[string]any{}, exitCode: exitSystemError, cause: err}
	}
	if errors.Is(err, io.EOF) && input == "" {
		return clientapp.Resource{}, pickCancelledError()
	}
	choice := strings.TrimSpace(input)
	if strings.EqualFold(choice, "q") {
		return clientapp.Resource{}, pickCancelledError()
	}
	if choice == "" {
		return clientapp.Resource{}, pickInvalidChoiceError(len(matches))
	}
	selected, convErr := strconv.Atoi(choice)
	if convErr != nil || selected < 1 || selected > len(matches) {
		return clientapp.Resource{}, pickInvalidChoiceError(len(matches))
	}
	return matches[selected-1], nil
}

func writePickValue(w io.Writer, value string) {
	value = strings.TrimRight(value, "\r\n")
	fmt.Fprintln(w, value)
}

func pickNotFoundError(query string) error {
	return &cliError{code: "not_found", message: i18n.T("command.pick.not_found", query), details: map[string]any{}, exitCode: exitNotFound}
}

func pickInvalidChoiceError(count int) error {
	return &cliError{code: "invalid_choice", message: i18n.T("command.pick.invalid_choice", count), details: map[string]any{"count": count}, exitCode: exitConflict}
}

func pickCancelledError() error {
	return &cliError{code: "cancelled", message: i18n.T("command.pick.cancelled"), details: map[string]any{}, exitCode: exitConflict}
}

func newUpdateCommand(_ *options) *cobra.Command {
	return &cobra.Command{
		Use:     "update <key> <new_value>",
		Aliases: []string{"put", "mod"},
		Short:   i18n.T("command.update.short"),
		Long:    i18n.T("command.update.long"),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := modeFromCommand(cmd)
			value, err := readValue(cmd, args[1])
			if err != nil {
				return err
			}
			if args[1] != "-" {
				value = util.UnescapeString(value)
			}
			debug, _ := cmd.Context().Value("debug").(bool)
			service := serviceFromCommand(cmd)
			if err := service.RecordAudit(args[0], "update"); err != nil && debug && !mode.json {
				fmt.Fprintf(cmd.OutOrStdout(), i18n.T("command.update.audit_error"), err)
			}
			resource, err := service.UpdateResourceValue(args[0], value)
			if err != nil {
				return textCommandError(cmd, "update", args[0], err)
			}
			if mode.json {
				return writeJSONSuccess(cmd.OutOrStdout(), resourceData{Resource: toResourceDTO(resource)})
			}
			return nil
		},
	}
}

func newDeleteCommand(_ *options) *cobra.Command {
	return &cobra.Command{
		Use:     "del [key]",
		Aliases: []string{"rm"},
		Short:   i18n.T("command.delete.short"),
		Long:    i18n.T("command.delete.long"),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := modeFromCommand(cmd)
			service := serviceFromCommand(cmd)
			result, err := service.DeleteResourceWithCleanup(args[0])
			if err != nil {
				if kind, ok := clientapp.ErrorKindOf(err); ok && kind == clientapp.ErrorNotFound && !mode.json && !mode.nonInteractive {
					fmt.Fprintf(cmd.OutOrStdout(), i18n.T("command.delete.not_found"), args[0])
					return nil
				}
				return textCommandError(cmd, "delete", args[0], err)
			}
			debug, _ := cmd.Context().Value("debug").(bool)
			if debug && !mode.json {
				if result.HistoryCleanupError != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "Failed to clean resource history: %v\n", result.HistoryCleanupError)
				}
				if result.AuditCleanupError != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "Failed to clean resource audit: %v\n", result.AuditCleanupError)
				}
			}
			if mode.json {
				return writeJSONSuccess(cmd.OutOrStdout(), deleteData{Key: args[0], Deleted: true})
			}
			fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.delete.success"))
			return nil
		},
	}
}

func newTagCommand(_ *options) *cobra.Command {
	return &cobra.Command{
		Use:   "tag [key] [tags...]",
		Short: i18n.T("command.tag.short"),
		Long:  i18n.T("command.tag.long"),
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			resource, err := serviceFromCommand(cmd).AddResourceTags(args[0], args[1:])
			if err != nil {
				return textCommandError(cmd, "tag", args[0], err)
			}
			if modeFromCommand(cmd).json {
				return writeJSONSuccess(cmd.OutOrStdout(), resourceData{Resource: toResourceDTO(resource)})
			}
			fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.tag.success"))
			return nil
		},
	}
}

func newDeleteTagCommand(_ *options) *cobra.Command {
	return &cobra.Command{
		Use:   "dtag [key] [tag]",
		Short: i18n.T("command.dtag.short"),
		Long:  i18n.T("command.dtag.long"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			resource, err := serviceFromCommand(cmd).DeleteResourceTag(args[0], args[1])
			if err != nil {
				return textCommandError(cmd, "dtag", args[0], err)
			}
			if modeFromCommand(cmd).json {
				return writeJSONSuccess(cmd.OutOrStdout(), resourceData{Resource: toResourceDTO(resource)})
			}
			fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.dtag.success"))
			return nil
		},
	}
}

func serviceFromCommand(cmd *cobra.Command) *clientapp.Service {
	service, err := clientapp.ServiceFromContext(cmd.Context())
	if err != nil {
		return clientapp.NewService(nil)
	}
	return service
}

type textError struct {
	message string
	cause   error
}

func (e *textError) Error() string { return e.message }
func (e *textError) Unwrap() error { return e.cause }

func localizedError(message string, cause error) error {
	return &textError{message: message, cause: cause}
}

func textCommandError(cmd *cobra.Command, commandName, key string, err error) error {
	if modeFromCommand(cmd).json {
		return err
	}
	var serviceErr *clientapp.ServiceError
	if !errors.As(err, &serviceErr) {
		return err
	}
	switch serviceErr.Kind {
	case clientapp.ErrorNotFound:
		switch commandName {
		case "get":
			return localizedError(i18n.T("command.get.not_found", key), err)
		case "update":
			return localizedError(i18n.T("command.update.not_found", key), err)
		case "tag":
			return localizedError(i18n.T("command.tag.not_found"), err)
		case "dtag":
			return localizedError(i18n.T("command.dtag.not_found"), err)
		}
	case clientapp.ErrorConflict:
		if commandName == "add" {
			return localizedError(i18n.T("command.add.duplicate", key), err)
		}
	case clientapp.ErrorSystem:
		return localizedSystemError(commandName, serviceErr)
	}
	return err
}

func localizedSystemError(commandName string, serviceErr *clientapp.ServiceError) error {
	var key string
	switch serviceErr.Operation {
	case clientapp.ErrorRead:
		key = "command." + commandName + ".error_fetch"
	case clientapp.ErrorSave:
		key = "command." + commandName + ".error_save"
	case clientapp.ErrorDelete:
		key = "command.delete.error_delete"
	}
	if key == "" {
		return serviceErr
	}
	cause := serviceErr.Err
	if cause == nil {
		cause = serviceErr
	}
	return fmt.Errorf(i18n.T(key), cause)
}

func selectResource(cmd *cobra.Command, matches []clientapp.Resource) (clientapp.Resource, error) {
	if len(matches) == 0 {
		return clientapp.Resource{}, interactionRequired("no resources available for selection")
	}
	matches = append([]clientapp.Resource(nil), matches...)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Key.Key < matches[j].Key.Key
	})
	fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.get.multiple_matches"))
	for index, match := range matches {
		fmt.Fprintf(cmd.OutOrStdout(), "%d. ", index+1)
		if match.Key.Type == models.ORIGIN {
			fmt.Fprintln(cmd.OutOrStdout(), match.Key.Key)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s "+i18n.T("command.get.tag_hint")+"\n", match.Key.OriginKey, match.Key.Key)
		}
	}
	var choice int
	if _, err := fmt.Fscan(cmd.InOrStdin(), &choice); err != nil {
		return clientapp.Resource{}, fmt.Errorf(i18n.T("command.get.invalid_input"), err)
	}
	if choice < 1 || choice > len(matches) {
		return clientapp.Resource{}, fmt.Errorf("%s", i18n.T("command.get.invalid_choice"))
	}
	return matches[choice-1], nil
}

func printTextResource(cmd *cobra.Command, resource clientapp.Resource) {
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.get.resource_label"), resourceKey(resource.Key))
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), resource.Value.Val)
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), i18n.T("command.get.tags_label"), resource.Value.Tag)
}
