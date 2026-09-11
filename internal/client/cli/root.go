package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"ttl-cli/command"
	"ttl-cli/conf"
	"ttl-cli/db"
	"ttl-cli/i18n"
	"ttl-cli/internal/client/remote"
	ttlsync "ttl-cli/sync"

	"github.com/spf13/cobra"
)

type options struct {
	debug        bool
	storageType  string
	cloudAPIURL  string
	cloudAPIKey  string
	cloudTimeout int
	confFile     string
}

// NewRootCommand builds the ttl client command tree.
func NewRootCommand() *cobra.Command {
	opts := &options{}
	root := &cobra.Command{
		Use:   "ttl",
		Short: i18n.T("root.short"),
		Long:  i18n.T("root.long"),
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	root.PersistentFlags().BoolVarP(&opts.debug, "debug", "D", false, i18n.T("root.flag_debug"))
	root.PersistentFlags().StringVar(&opts.storageType, "storage", "sqlite", i18n.T("root.flag_storage"))
	root.PersistentFlags().StringVar(&opts.cloudAPIURL, "cloud-url", "", i18n.T("root.flag_cloud_url"))
	root.PersistentFlags().StringVar(&opts.cloudAPIKey, "cloud-key", "", i18n.T("root.flag_cloud_key"))
	root.PersistentFlags().IntVar(&opts.cloudTimeout, "cloud-timeout", 30, i18n.T("root.flag_cloud_timeout"))
	root.PersistentFlags().StringVar(&opts.confFile, "conf", "", i18n.T("root.flag_conf"))

	root.AddCommand(
		command.InitCmd,
		command.AddCmd,
		command.GetCmd,
		command.OpenCmd,
		command.UpdateCmd,
		command.DelCmd,
		command.TagCmd,
		command.DtagCmd,
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
		newServerCompatibilityCommand(),
		newSyncCommand(opts),
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

	root := NewRootCommand()
	updateCommandDescriptions(root)
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		return 1
	}

	debug, _ := root.PersistentFlags().GetBool("debug")
	if err := db.CloseDB(); err != nil && debug {
		fmt.Printf(i18n.T("error.close_db"), err)
	}
	return 0
}

func newPreRun(opts *options) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "server" || cmd.Parent() != nil && cmd.Parent().Name() == "user" {
			return nil
		}

		skipDBInit := cmd.Name() == "workspace" ||
			cmd.Parent() != nil && cmd.Parent().Name() == "workspace" ||
			cmd.Name() == "ws"

		ctx := context.WithValue(cmd.Context(), "debug", opts.debug)
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

			if err := db.InitDB(actualStorageType, opts.cloudAPIURL, opts.cloudAPIKey, opts.cloudTimeout, opts.confFile); err != nil {
				return fmt.Errorf(i18n.T("error.init_db"), err)
			}
			replaceSpecialValuesFromHistory(args)
			if shouldRecordHistory(cmd) {
				resourceKey := ""
				if len(args) > 0 {
					resourceKey = args[0]
				}
				if err := db.RecordCommandHistory(cmd.Name(), resourceKey, opts.debug); err != nil && opts.debug {
					fmt.Printf(i18n.T("error.record_history"), err)
				}
			}
		}
		cmd.SetContext(ctx)
		return nil
	}
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
			localResources, err := db.GetAllResources()
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
				return ttlsync.ExecutePull(diff, db.Stor, remoteStorage, false)
			case "push":
				return ttlsync.ExecutePush(diff, db.Stor, remoteStorage, false)
			case "auto":
				return executeInteractiveSync(diff, remoteStorage)
			default:
				return fmt.Errorf(i18n.T("command.sync.invalid_direction"), direction)
			}
		},
	}
	cmd.Flags().StringVar(&direction, "direction", "auto", i18n.T("command.sync.flag_direction"))
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, i18n.T("command.sync.flag_dry_run"))
	return cmd
}

func executeInteractiveSync(diff ttlsync.DiffResult, remoteStorage db.Storage) error {
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
		return ttlsync.ExecutePull(diff, db.Stor, remoteStorage, false)
	case "push":
		return ttlsync.ExecutePush(diff, db.Stor, remoteStorage, false)
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
			return db.MigrateData(sourceType, targetType, sourceAPIURL, sourceAPIKey, sourceTimeout,
				opts.cloudAPIURL, opts.cloudAPIKey, opts.cloudTimeout, opts.debug, opts.confFile, opts.confFile)
		},
	}
	cmd.Flags().String("source-url", "", i18n.T("command.migrate.flag_source_url"))
	cmd.Flags().String("source-key", "", i18n.T("command.migrate.flag_source_key"))
	cmd.Flags().Int("source-timeout", 30, i18n.T("command.migrate.flag_source_timeout"))
	return cmd
}

func replaceSpecialValuesFromHistory(args []string) {
	if len(args) != 1 {
		return
	}
	charsCount := countSpecialChars(args[0])
	if charsCount > 0 {
		record, err := db.GetHistoryRecords(charsCount - 1)
		if err != nil {
			fmt.Printf(i18n.T("error.get_history"), err)
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
