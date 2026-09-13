package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"ttl-cli/command"
	"ttl-cli/conf"
	"ttl-cli/i18n"
	clientapp "ttl-cli/internal/client/app"
	"ttl-cli/internal/client/remote"
	clienttui "ttl-cli/internal/client/tui"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"
	ttlsync "ttl-cli/sync"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

type storageOpener func(string, string, string, int, string) (corestorage.Storage, error)
type terminalDetector func(io.Reader, io.Writer) bool
type tuiRunner func(clienttui.ResourceService, clienttui.RunOptions) error

type options struct {
	debug          bool
	json           bool
	nonInteractive bool
	storageType    string
	cloudAPIURL    string
	cloudAPIKey    string
	cloudTimeout   int
	confFile       string
	service        *clientapp.Service
	openStorage    storageOpener
	isTerminal     terminalDetector
	runTUI         tuiRunner
}

type runResult struct {
	exitCode int
	stdout   []byte
}

// NewRootCommand builds the ttl client command tree.
func NewRootCommand() *cobra.Command {
	return newRootCommand(&options{})
}

func newRootCommand(opts *options) *cobra.Command {
	root := &cobra.Command{
		Use:   "ttl",
		Short: i18n.T("root.short"),
		Long:  i18n.T("root.long"),
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	root.PersistentFlags().BoolVarP(&opts.debug, "debug", "D", false, i18n.T("root.flag_debug"))
	root.PersistentFlags().BoolVar(&opts.json, "json", false, "Output versioned JSON for supported resource commands")
	root.PersistentFlags().BoolVar(&opts.nonInteractive, "non-interactive", false, "Disable interactive input for supported resource commands")
	root.PersistentFlags().StringVar(&opts.storageType, "storage", "sqlite", i18n.T("root.flag_storage"))
	root.PersistentFlags().StringVar(&opts.cloudAPIURL, "cloud-url", "", i18n.T("root.flag_cloud_url"))
	root.PersistentFlags().StringVar(&opts.cloudAPIKey, "cloud-key", "", i18n.T("root.flag_cloud_key"))
	root.PersistentFlags().IntVar(&opts.cloudTimeout, "cloud-timeout", 30, i18n.T("root.flag_cloud_timeout"))
	root.PersistentFlags().StringVar(&opts.confFile, "conf", "", i18n.T("root.flag_conf"))

	root.AddCommand(
		command.InitCmd,
		newAddCommand(opts),
		newGetCommand(opts),
		command.OpenCmd,
		newUpdateCommand(opts),
		newDeleteCommand(opts),
		newTagCommand(opts),
		newDeleteTagCommand(opts),
		command.TagsCmd,
		command.RenameCmd,
		command.ConfigCmd,
		command.VersionCmd,
		command.EncryptCmd,
		command.DecryptCmd,
		command.KeyCmd,
		newMigrateCommand(opts),
		command.AuditCmd,
		command.HistoryCmd,
		command.ExportCmd,
		command.ImportCmd,
		command.LogCmd,
		newSyncCommand(opts),
		newUICommand(opts),
		command.WorkspaceCmd,
		command.WsCmd,
	)

	root.PersistentPreRunE = newPreRun(opts)
	return root
}

// Run initializes localization, executes the ttl client, and closes storage.
func Run() int {
	if err := i18n.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize i18n: %v\n", err)
	}

	opts := &options{}
	root := newRootCommand(opts)
	root.SetIn(os.Stdin)
	root.SetErr(os.Stderr)
	requestedJSON := boolFlagEnabled(os.Args[1:], "--json")
	requestedNonInteractive := boolFlagEnabled(os.Args[1:], "--non-interactive")
	if requestedJSON || requestedNonInteractive {
		root.SilenceErrors = true
		root.SilenceUsage = true
	}
	var jsonOutput bytes.Buffer
	if requestedJSON {
		root.SetOut(&jsonOutput)
	} else {
		root.SetOut(os.Stdout)
	}
	updateCommandDescriptions(root)
	result := executeRoot(root, opts, os.Args[1:], requestedJSON, &jsonOutput)
	if len(result.stdout) > 0 {
		if _, err := os.Stdout.Write(result.stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitSystemError
		}
	}
	return result.exitCode
}

func executeRoot(root *cobra.Command, opts *options, args []string, requestedJSON bool, jsonOutput *bytes.Buffer) (result runResult) {
	defer func() {
		if requestedJSON && result.exitCode == exitSuccess {
			result.stdout = append([]byte(nil), jsonOutput.Bytes()...)
		}
	}()
	root.SetArgs(args)
	executeErr := root.Execute()
	requestedNonInteractive := boolFlagEnabled(args, "--non-interactive")
	if executeErr != nil && (requestedJSON || requestedNonInteractive) {
		if _, ok := executeErr.(*cliError); !ok && isArgumentError(executeErr) {
			executeErr = invalidArgument(executeErr.Error())
		}
	}

	executeErr = mergeCloseError(root, opts, executeErr)
	if executeErr != nil {
		jsonMode := opts.json || requestedJSON
		machineMode := jsonMode || opts.nonInteractive || requestedNonInteractive
		if jsonMode {
			if err := writeJSONError(root.ErrOrStderr(), executeErr); err != nil {
				fmt.Fprintln(root.ErrOrStderr(), err)
				return runResult{exitCode: exitSystemError}
			}
		} else if machineMode {
			fmt.Fprintln(root.ErrOrStderr(), executeErr)
		} else {
			fmt.Fprintln(root.OutOrStdout(), executeErr)
		}
		return runResult{exitCode: exitCodeFor(executeErr, machineMode)}
	}
	return runResult{exitCode: exitSuccess}
}

func mergeCloseError(root *cobra.Command, opts *options, executeErr error) error {
	if opts.service == nil {
		return executeErr
	}
	debug, _ := root.PersistentFlags().GetBool("debug")
	closeErr := opts.service.Close()
	if executeErr == nil && closeErr != nil {
		return closeErr
	}
	if closeErr != nil && debug && !opts.json {
		fmt.Fprintf(root.OutOrStdout(), i18n.T("error.close_db"), closeErr)
	}
	return executeErr
}

func boolFlagEnabled(args []string, target string) bool {
	enabled := false
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == target {
			enabled = true
			continue
		}
		prefix := target + "="
		if strings.HasPrefix(arg, prefix) {
			value, err := strconv.ParseBool(strings.TrimPrefix(arg, prefix))
			if err != nil {
				enabled = true
				continue
			}
			enabled = value
		}
	}
	return enabled
}

func newPreRun(opts *options) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		mode := invocationMode{json: opts.json, nonInteractive: opts.nonInteractive || opts.json}
		ctx := context.WithValue(cmd.Context(), machineModeKey{}, mode)
		cmd.SetContext(ctx)
		if err := validateMachineMode(cmd); err != nil {
			return err
		}
		if cmd.Name() == "server" || cmd.Parent() != nil && cmd.Parent().Name() == "user" {
			return nil
		}

		skipDBInit := cmd.Name() == "workspace" ||
			cmd.Parent() != nil && cmd.Parent().Name() == "workspace" ||
			cmd.Name() == "ws"

		ctx = context.WithValue(cmd.Context(), "debug", opts.debug)
		ctx = context.WithValue(ctx, "confFile", opts.confFile)
		if !skipDBInit {
			actualStorageType := opts.storageType
			if opts.storageType == "sqlite" && !cmd.Flags().Changed("storage") {
				ttlConf, err := conf.GetTtlConfFromFile(opts.confFile)
				if err == nil {
					if ttlConf.Workspace != "" {
						if ws, ok := ttlConf.Workspaces[ttlConf.Workspace]; ok && ws.StorageType != "" {
							actualStorageType = ws.StorageType
						}
					}
					if actualStorageType == "sqlite" && ttlConf.StorageType != "" {
						actualStorageType = ttlConf.StorageType
					}
				}
			}
			if cmd.Name() == "ui" {
				switch actualStorageType {
				case "sqlite", "local", "bbolt":
				default:
					return fmt.Errorf("ttl ui only supports local sqlite or bbolt storage")
				}
				detector := opts.isTerminal
				if detector == nil {
					detector = defaultTerminalDetector
				}
				if !detector(cmd.InOrStdin(), cmd.OutOrStdout()) {
					return fmt.Errorf("ttl ui requires an interactive terminal; use the CLI in non-interactive environments")
				}
			}
			if opts.service != nil {
				_ = opts.service.Close()
			}

			opener := opts.openStorage
			if opener == nil {
				opener = clientapp.OpenStorage
			}
			storage, err := opener(actualStorageType, opts.cloudAPIURL, opts.cloudAPIKey, opts.cloudTimeout, opts.confFile)
			if err != nil {
				return fmt.Errorf(i18n.T("error.init_db"), err)
			}
			opts.service = clientapp.NewService(storage)
			replaceSpecialValuesFromHistory(cmd, opts.service, args)
			if shouldRecordHistory(cmd) {
				resourceKey := ""
				if len(args) > 0 {
					resourceKey = args[0]
				}
				if err := opts.service.RecordCommandHistory(cmd.Name(), resourceKey, opts.debug); err != nil && opts.debug && !mode.json {
					fmt.Fprintf(cmd.OutOrStdout(), i18n.T("error.record_history"), err)
				}
			}
		}
		cmd.SetContext(clientapp.WithService(ctx, opts.service))
		return nil
	}
}

func newUICommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "ui",
		Short: i18n.T("command.ui.short"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			runner := opts.runTUI
			if runner == nil {
				runner = clienttui.Run
			}
			return runner(serviceFromCommand(cmd), clienttui.RunOptions{In: cmd.InOrStdin(), Out: cmd.OutOrStdout()})
		},
	}
}

func defaultTerminalDetector(input io.Reader, output io.Writer) bool {
	in, inOK := input.(*os.File)
	out, outOK := output.(*os.File)
	return inOK && outOK && term.IsTerminal(in.Fd()) && term.IsTerminal(out.Fd())
}

func shouldRecordHistory(cmd *cobra.Command) bool {
	if cmd.HasSubCommands() {
		return false
	}
	switch cmd.Name() {
	case "history", "audit", "export", "server", "sync", "log", "tags":
		return false
	default:
		return true
	}
}

func newSyncCommand(opts *options) *cobra.Command {
	var direction string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "sync",
		Short: i18n.T("command.sync.short"),
		Long:  i18n.T("command.sync.long"),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if opts.cloudAPIURL == "" {
				return errors.New(i18n.T("command.sync.need_cloud_url"))
			}
			localResources, err := opts.service.GetAllResources()
			if err != nil {
				return fmt.Errorf(i18n.T("command.sync.error_fetch_local"), err)
			}
			remoteStorage := remote.NewStorage(opts.cloudAPIURL, opts.cloudAPIKey, opts.cloudTimeout)
			if err := remoteStorage.Init(); err != nil {
				return fmt.Errorf(i18n.T("command.sync.error_connect_remote"), err)
			}
			defer remoteStorage.Close()
			remoteResources, err := remoteStorage.GetAllResources()
			if err != nil {
				return fmt.Errorf(i18n.T("command.sync.error_fetch_remote"), err)
			}

			diff := ttlsync.ComputeDiff(localResources, remoteResources)
			ttlsync.PrintDiff(diff, opts.cloudAPIURL)
			if diff.InSync || dryRun {
				if dryRun && !diff.InSync {
					fmt.Println(i18n.T("command.sync.dry_run_notice"))
				}
				return nil
			}

			switch direction {
			case "pull":
				return ttlsync.ExecutePull(diff, opts.service.Storage(), remoteStorage, false)
			case "push":
				return ttlsync.ExecutePush(diff, opts.service.Storage(), remoteStorage, false)
			case "auto":
				return executeInteractiveSync(diff, opts.service.Storage(), remoteStorage)
			default:
				return fmt.Errorf(i18n.T("command.sync.invalid_direction"), direction)
			}
		},
	}
	cmd.Flags().StringVar(&direction, "direction", "auto", i18n.T("command.sync.flag_direction"))
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, i18n.T("command.sync.flag_dry_run"))
	return cmd
}

func executeInteractiveSync(diff ttlsync.DiffResult, localStorage, remoteStorage corestorage.Storage) error {
	fmt.Println(i18n.T("command.sync.choose_operation"))
	fmt.Println(i18n.T("command.sync.option_pull"))
	fmt.Println(i18n.T("command.sync.option_push"))
	fmt.Println(i18n.T("command.sync.option_skip"))
	fmt.Print("> ")
	var choice string
	if _, err := fmt.Scan(&choice); err != nil {
		return fmt.Errorf(i18n.T("command.sync.invalid_input"), err)
	}
	switch choice {
	case "pull":
		return ttlsync.ExecutePull(diff, localStorage, remoteStorage, false)
	case "push":
		return ttlsync.ExecutePush(diff, localStorage, remoteStorage, false)
	case "skip":
		fmt.Println(i18n.T("command.sync.skipped"))
		return nil
	default:
		return fmt.Errorf(i18n.T("command.sync.invalid_choice"), choice)
	}
}

func newMigrateCommand(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [source] [target]",
		Short: i18n.T("command.migrate.short"),
		Long:  i18n.T("command.migrate.long"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceType, targetType := args[0], args[1]
			if sourceType != "local" && sourceType != "cloud" {
				return errors.New(i18n.T("command.migrate.invalid_source"))
			}
			if targetType != "local" && targetType != "cloud" {
				return errors.New(i18n.T("command.migrate.invalid_target"))
			}
			if sourceType == targetType {
				return errors.New(i18n.T("command.migrate.same_type"))
			}

			var sourceAPIURL, sourceAPIKey string
			var sourceTimeout int
			if sourceType == "cloud" {
				sourceAPIURL, _ = cmd.Flags().GetString("source-url")
				sourceAPIKey, _ = cmd.Flags().GetString("source-key")
				sourceTimeout, _ = cmd.Flags().GetInt("source-timeout")
				if sourceAPIURL == "" || sourceAPIKey == "" {
					return errors.New(i18n.T("command.migrate.need_source_config"))
				}
			}
			return clientapp.MigrateData(sourceType, targetType, sourceAPIURL, sourceAPIKey, sourceTimeout,
				opts.cloudAPIURL, opts.cloudAPIKey, opts.cloudTimeout, opts.debug, opts.confFile, opts.confFile)
		},
	}
	cmd.Flags().String("source-url", "", i18n.T("command.migrate.flag_source_url"))
	cmd.Flags().String("source-key", "", i18n.T("command.migrate.flag_source_key"))
	cmd.Flags().Int("source-timeout", 30, i18n.T("command.migrate.flag_source_timeout"))
	return cmd
}

func replaceSpecialValuesFromHistory(cmd *cobra.Command, service *clientapp.Service, args []string) {
	if len(args) != 1 {
		return
	}
	charsCount := countSpecialChars(args[0])
	if charsCount > 0 {
		record, err := service.GetHistoryRecord(charsCount-1, models.Descending)
		if err != nil {
			if !modeFromCommand(cmd).json {
				fmt.Fprintf(cmd.OutOrStdout(), i18n.T("error.get_history"), err)
			}
			return
		}
		args[0] = record.ResourceKey
	}
}

func countSpecialChars(input string) int {
	count := 0
	for _, char := range input {
		if char != '~' && char != '^' {
			return 0
		}
		count++
	}
	return count
}

func updateCommandDescriptions(cmd *cobra.Command) {
	if cmd.Short != "" {
		cmd.Short = i18n.T(cmd.Short)
	}
	if cmd.Long != "" {
		cmd.Long = i18n.T(cmd.Long)
	}
	for _, subCmd := range cmd.Commands() {
		updateCommandDescriptions(subCmd)
	}
}
