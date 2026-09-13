package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ttl-cli/internal/i18n"
	api "ttl-cli/internal/server/api"
	"ttl-cli/internal/server/tenant"

	"github.com/spf13/cobra"
)

type commandConfig struct {
	port    int
	dataDir string
}

// NewCompatibilityCommand returns the legacy server command tree for callers
// that still need to embed the standalone server command.
func NewCompatibilityCommand() *cobra.Command {
	cfg := defaultConfig()
	cmd := newServeCommand("server", &cfg)
	cmd.AddCommand(newUserCommand(&cfg))
	return cmd
}

// NewRootCommand returns the standalone ttl-server command tree.
func NewRootCommand() *cobra.Command {
	cfg := defaultConfig()
	root := &cobra.Command{
		Use:   "ttl-server",
		Short: i18n.T("command.server.short"),
		Long:  i18n.T("command.server.long"),
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}
	root.PersistentFlags().IntVar(&cfg.port, "port", 8080, i18n.T("command.server.flag_port"))
	root.PersistentFlags().StringVar(&cfg.dataDir, "data-dir", cfg.dataDir, i18n.T("command.server.flag_data_dir"))

	root.AddCommand(newServeCommand("serve", &cfg))
	root.AddCommand(newUserCommand(&cfg))
	return root
}

// Run initializes localization and executes the standalone server command.
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
	return 0
}

func defaultConfig() commandConfig {
	home, _ := os.UserHomeDir()
	return commandConfig{
		port:    8080,
		dataDir: filepath.Join(home, ".ttl"),
	}
}

func newServeCommand(use string, cfg *commandConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: i18n.T("command.server.short"),
		Long:  i18n.T("command.server.long"),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return api.StartServer(cfg.port, cfg.dataDir)
		},
	}
	if use == "server" {
		cmd.Flags().IntVar(&cfg.port, "port", 8080, i18n.T("command.server.flag_port"))
		cmd.Flags().StringVar(&cfg.dataDir, "data-dir", cfg.dataDir, i18n.T("command.server.flag_data_dir"))
	}
	return cmd
}

func newUserCommand(cfg *commandConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: i18n.T("command.server.user.short"),
		Long:  i18n.T("command.server.user.long"),
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(
		newUserAddCommand(cfg),
		newUserListCommand(cfg),
		newUserActiveCommand(cfg, true),
		newUserActiveCommand(cfg, false),
		newUserResetKeyCommand(cfg),
		newUserDeleteCommand(cfg),
	)
	return cmd
}

func newUserAddCommand(cfg *commandConfig) *cobra.Command {
	var id, name string
	cmd := &cobra.Command{
		Use:   "add",
		Short: i18n.T("command.server.user_add.short"),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			store, err := loadUserStore(cfg.dataDir)
			if err != nil {
				return err
			}
			user, err := store.AddUser(id, name)
			if err != nil {
				return err
			}
			fmt.Println(i18n.T("command.server.user_add.success"))
			fmt.Printf("  ID:      %s\n", user.ID)
			fmt.Printf("  Name:    %s\n", user.Name)
			fmt.Printf("  API Key: %s\n", user.APIKey)
			fmt.Println(i18n.T("command.server.user_add.warn_keep_key"))
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", i18n.T("command.server.user_add.flag_id"))
	cmd.Flags().StringVar(&name, "name", "", i18n.T("command.server.user_add.flag_name"))
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newUserListCommand(cfg *commandConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: i18n.T("command.server.user_list.short"),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			store, err := loadUserStore(cfg.dataDir)
			if err != nil {
				return err
			}
			users := store.ListUsers()
			if len(users) == 0 {
				fmt.Println(i18n.T("command.server.user_list.no_users"))
				return nil
			}
			fmt.Printf("%-16s %-16s %-8s %-12s %s\n", "ID", "Name", "Active", "API Key", "Created")
			for _, user := range users {
				keyPreview := user.APIKey
				if len(keyPreview) > 8 {
					keyPreview = keyPreview[:8] + "****"
				}
				created := time.Unix(user.CreatedAt, 0).Format("2006-01-02 15:04:05")
				fmt.Printf("%-16s %-16s %-8v %-12s %s\n", user.ID, user.Name, user.Active, keyPreview, created)
			}
			return nil
		},
	}
}

func newUserActiveCommand(cfg *commandConfig, active bool) *cobra.Command {
	name := "disable"
	shortKey := "command.server.user_disable.short"
	successKey := "command.server.user_disable.success"
	flagKey := "command.server.user_disable.flag_id"
	if active {
		name = "enable"
		shortKey = "command.server.user_enable.short"
		successKey = "command.server.user_enable.success"
		flagKey = "command.server.user_enable.flag_id"
	}

	var id string
	cmd := &cobra.Command{
		Use:   name,
		Short: i18n.T(shortKey),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			store, err := loadUserStore(cfg.dataDir)
			if err != nil {
				return err
			}
			if err := store.SetActive(id, active); err != nil {
				return err
			}
			fmt.Printf(i18n.T(successKey), id)
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", i18n.T(flagKey))
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newUserResetKeyCommand(cfg *commandConfig) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "reset-key",
		Short: i18n.T("command.server.user_reset_key.short"),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			store, err := loadUserStore(cfg.dataDir)
			if err != nil {
				return err
			}
			newKey, err := store.ResetKey(id)
			if err != nil {
				return err
			}
			fmt.Printf(i18n.T("command.server.user_reset_key.success"), id)
			fmt.Printf(i18n.T("command.server.user_reset_key.new_key"), newKey)
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", i18n.T("command.server.user_reset_key.flag_id"))
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newUserDeleteCommand(cfg *commandConfig) *cobra.Command {
	var id string
	var confirm bool
	cmd := &cobra.Command{
		Use:   "delete",
		Short: i18n.T("command.server.user_delete.short"),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if !confirm {
				return errors.New(i18n.T("command.server.user_delete.need_confirm"))
			}
			store, err := loadUserStore(cfg.dataDir)
			if err != nil {
				return err
			}
			if err := store.DeleteUser(id); err != nil {
				return err
			}

			tenantMgr := tenant.NewStorageManager(filepath.Join(cfg.dataDir, "tenants"))
			if err := tenantMgr.RemoveStorage(id); err != nil {
				fmt.Printf(i18n.T("command.server.user_delete.warn_delete_data"), err)
			}
			fmt.Printf(i18n.T("command.server.user_delete.success"), id)
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", i18n.T("command.server.user_delete.flag_id"))
	cmd.Flags().BoolVar(&confirm, "confirm", false, i18n.T("command.server.user_delete.flag_confirm"))
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func loadUserStore(dataDir string) (*tenant.UserStore, error) {
	store := tenant.NewUserStore(filepath.Join(dataDir, "users.json"))
	if err := store.Load(); err != nil {
		return nil, err
	}
	return store, nil
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
