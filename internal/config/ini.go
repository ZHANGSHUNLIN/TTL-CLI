package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const DefaultStorageType = "local"

// ValidateStorageType rejects removed storage backends instead of silently selecting one.
func ValidateStorageType(storageType string) error {
	switch storageType {
	case "local", "cloud":
		return nil
	case "":
		return nil
	case "sync":
		return fmt.Errorf("不支持的存储模式 %q；sync 是独立同步能力，请使用 ttl sync", storageType)
	default:
		return fmt.Errorf("不支持的存储模式 %q，仅支持 local 或 cloud", storageType)
	}
}

func GetTtlConf() (TtlIni, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return TtlIni{}, fmt.Errorf("failed to get user directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".ttl")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return TtlIni{}, fmt.Errorf("failed to create config directory: %w", err)
	}

	confFilePath := filepath.Join(configDir, "ttl.ini")

	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		return createDefaultConfig(confFilePath, "")
	}

	return loadConfFile(confFilePath)
}

func GetTtlConfFromFile(confFile string) (TtlIni, error) {
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(confFile), 0755); err != nil {
			return TtlIni{}, fmt.Errorf("failed to create config directory: %w", err)
		}
		return createDefaultConfig(confFile, "")
	}
	return loadConfFile(confFile)
}

func loadConfFile(path string) (TtlIni, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return TtlIni{}, fmt.Errorf("failed to load config file: %w", err)
	}
	var ttlIni TtlIni
	if err := cfg.Section("").MapTo(&ttlIni); err != nil {
		return TtlIni{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.HasSection("storage") {
		storageSec := cfg.Section("storage")
		if storageType := storageSec.Key("type").String(); storageType != "" {
			ttlIni.StorageType = storageType
		}
		if storagePath := storageSec.Key("path").String(); storagePath != "" {
			ttlIni.DbPath = storagePath
		}
		if remoteName := storageSec.Key("remote").String(); remoteName != "" {
			ttlIni.RemoteName = remoteName
		}
	}

	if ttlIni.StorageType == "" {
		ttlIni.StorageType = DefaultStorageType
	}
	if err := ValidateStorageType(ttlIni.StorageType); err != nil {
		return TtlIni{}, err
	}

	if err := cfg.Section("bbolt").MapTo(&ttlIni.BoltDB); err != nil {
		return TtlIni{}, fmt.Errorf("failed to parse bbolt config: %w", err)
	}
	if ttlIni.BoltDB.Timeout == 0 {
		ttlIni.BoltDB.Timeout = 5
	}

	ttlIni.Workspaces = make(map[string]WorkspaceConfig)
	ttlIni.Remotes = make(map[string]RemoteConfig)
	for _, section := range cfg.Sections() {
		name := section.Name()
		if strings.HasPrefix(name, "workspaces.") {
			wsName := strings.TrimPrefix(name, "workspaces.")
			if wsName != "" {
				wsConfig := WorkspaceConfig{
					DbPath:      section.Key("db_path").String(),
					StorageType: section.Key("storage_type").String(),
					RemoteName:  section.Key("remote").String(),
				}
				if err := ValidateStorageType(wsConfig.StorageType); err != nil {
					return TtlIni{}, fmt.Errorf("工作空间 %s: %w", wsName, err)
				}
				if wsConfig.StorageType == "" {
					wsConfig.StorageType = ttlIni.StorageType
				}
				ttlIni.Workspaces[wsName] = wsConfig
			}
		}
		if strings.HasPrefix(name, "remotes.") {
			remoteName := strings.TrimPrefix(name, "remotes.")
			if remoteName != "" {
				ttlIni.Remotes[remoteName] = RemoteConfig{
					URL:           section.Key("url").String(),
					Account:       section.Key("account").String(),
					CredentialEnv: section.Key("credential_env").String(),
				}
			}
		}
	}

	return ttlIni, nil
}

func GetDefaultConfPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user directory: %w", err)
	}
	return filepath.Join(homeDir, ".ttl", "ttl.ini"), nil
}

func createDefaultConfig(configPath, dbPath string) (TtlIni, error) {
	cfg := ini.Empty()

	storageSec := cfg.Section("storage")
	storageSec.Key("type").SetValue(DefaultStorageType)

	if dbPath != "" {
		storageSec.Key("path").SetValue(dbPath)
		cfg.Section("").Key("db_path").SetValue(dbPath)
	}

	if err := cfg.SaveTo(configPath); err != nil {
		return TtlIni{}, fmt.Errorf("failed to save default config: %w", err)
	}

	return loadConfFile(configPath)
}

func GetWorkspaceDBPath(confFile string) (string, string, error) {
	var ttlConf TtlIni
	var err error
	if confFile != "" {
		ttlConf, err = GetTtlConfFromFile(confFile)
	} else {
		ttlConf, err = GetTtlConf()
	}
	if err != nil {
		return "", "", err
	}

	workspaceName := ttlConf.Workspace
	if workspaceName == "" {
		workspaceName = "default"
	}

	if ws, ok := ttlConf.Workspaces[workspaceName]; ok && ws.DbPath != "" {
		return workspaceName, ws.DbPath, nil
	}

	if ttlConf.DbPath != "" {
		return workspaceName, ttlConf.DbPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to get user directory: %w", err)
	}
	return workspaceName, filepath.Join(homeDir, ".ttl", "data.sqlite"), nil
}

// LoadRemote returns the active remote profile for the current workspace.
func LoadRemote(confFile string) (string, RemoteConfig, error) {
	var ttlConf TtlIni
	var err error
	if confFile != "" {
		ttlConf, err = GetTtlConfFromFile(confFile)
	} else {
		ttlConf, err = GetTtlConf()
	}
	if err != nil {
		return "", RemoteConfig{}, err
	}
	remoteName := ttlConf.RemoteName
	if ttlConf.Workspace != "" {
		if ws, ok := ttlConf.Workspaces[ttlConf.Workspace]; ok && ws.RemoteName != "" {
			remoteName = ws.RemoteName
		}
	}
	if remoteName == "" {
		return "", RemoteConfig{}, fmt.Errorf("未配置活动远程 profile")
	}
	remote, ok := ttlConf.Remotes[remoteName]
	if !ok || remote.URL == "" {
		return "", RemoteConfig{}, fmt.Errorf("远程 profile %q 未配置地址", remoteName)
	}
	return remoteName, remote, nil
}

func ValidateWorkspaceName(name string) bool {
	if name == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched
}

func CreateWorkspace(confFile, name string) (string, error) {
	if !ValidateWorkspaceName(name) {
		return "", fmt.Errorf("工作空间名称只能包含字母、数字、下划线、连字符")
	}

	path := confFile
	if path == "" {
		var err error
		path, err = GetDefaultConfPath()
		if err != nil {
			return "", err
		}
	}

	cfg, err := ini.Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = ini.Empty()
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return "", fmt.Errorf("failed to create config directory: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to load config file: %w", err)
		}
	}

	sectionName := "workspaces." + name
	if cfg.HasSection(sectionName) {
		return "", fmt.Errorf("工作空间已存在: %s", name)
	}

	// 确定工作空间目录：基于配置文件所在目录
	var wsDir string
	confDir := filepath.Dir(path)
	// 如果配置文件在 ~/.ttl/ 目录下，则使用 ~/.ttl/workspaces/
	homeDir, err := os.UserHomeDir()
	if err == nil {
		defaultConfDir := filepath.Join(homeDir, ".ttl")
		if confDir == defaultConfDir {
			wsDir = filepath.Join(homeDir, ".ttl", "workspaces")
		}
	}
	// 否则在配置文件同目录下创建 workspaces 子目录
	if wsDir == "" {
		wsDir = filepath.Join(confDir, "workspaces")
	}
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create workspaces directory: %w", err)
	}

	storageType := ""
	if cfg.HasSection("storage") {
		storageType = cfg.Section("storage").Key("type").String()
		if storageType == "" {
			storageType = cfg.Section("storage").Key("storage_type").String()
		}
	}
	if storageType == "" {
		storageType = cfg.Section("").Key("storage_type").String()
	}
	if storageType == "" {
		storageType = DefaultStorageType
	}

	if err := ValidateStorageType(storageType); err != nil {
		return "", err
	}

	dbPath := filepath.Join(wsDir, name+".sqlite")

	sec := cfg.Section(sectionName)
	sec.Key("db_path").SetValue(dbPath)
	sec.Key("storage_type").SetValue(storageType)

	if err := cfg.SaveTo(path); err != nil {
		return "", fmt.Errorf("failed to save config: %w", err)
	}

	return dbPath, nil
}

func DeleteWorkspace(confFile, name string) error {
	path := confFile
	if path == "" {
		var err error
		path, err = GetDefaultConfPath()
		if err != nil {
			return err
		}
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	sectionName := "workspaces." + name
	if !cfg.HasSection(sectionName) {
		return fmt.Errorf("工作空间不存在: %s", name)
	}

	sec := cfg.Section(sectionName)
	dbPath := sec.Key("db_path").String()

	cfg.DeleteSection(sectionName)

	if err := cfg.SaveTo(path); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if dbPath != "" {
		if _, err := os.Stat(dbPath); err == nil {
			if err := os.Remove(dbPath); err != nil {
				return fmt.Errorf("failed to delete database file: %w", err)
			}
		}
	}

	keyFile := dbPath[:len(dbPath)-len(filepath.Ext(dbPath))] + ".key"
	if _, err := os.Stat(keyFile); err == nil {
		os.Remove(keyFile)
	}

	return nil
}

func SwitchWorkspace(confFile, name string) error {
	path := confFile
	if path == "" {
		var err error
		path, err = GetDefaultConfPath()
		if err != nil {
			return err
		}
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	sectionName := "workspaces." + name
	if !cfg.HasSection(sectionName) && name != "default" {
		return fmt.Errorf("工作空间不存在: %s", name)
	}

	workspaceSec := cfg.Section("")
	workspaceSec.Key("workspace").SetValue(name)

	if err := cfg.SaveTo(path); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func GetCurrentWorkspace(confFile string) (string, error) {
	var ttlConf TtlIni
	var err error
	if confFile != "" {
		ttlConf, err = GetTtlConfFromFile(confFile)
	} else {
		ttlConf, err = GetTtlConf()
	}
	if err != nil {
		return "", err
	}

	if ttlConf.Workspace == "" {
		return "default", nil
	}
	return ttlConf.Workspace, nil
}

func ListWorkspaces(confFile string) ([]string, string, error) {
	var ttlConf TtlIni
	var err error
	if confFile != "" {
		ttlConf, err = GetTtlConfFromFile(confFile)
	} else {
		ttlConf, err = GetTtlConf()
	}
	if err != nil {
		return nil, "", err
	}

	current := ttlConf.Workspace
	if current == "" {
		current = "default"
	}

	var names []string
	for name := range ttlConf.Workspaces {
		names = append(names, name)
	}

	return names, current, nil
}

func GetWorkspaceInfo(confFile, name string) (string, string, int, error) {
	var ttlConf TtlIni
	var err error
	if confFile != "" {
		ttlConf, err = GetTtlConfFromFile(confFile)
	} else {
		ttlConf, err = GetTtlConf()
	}
	if err != nil {
		return "", "", 0, err
	}

	var dbPath, storageType string

	if name == "default" {
		dbPath = ttlConf.DbPath
		storageType = ttlConf.StorageType
		if storageType == "" {
			storageType = DefaultStorageType
		}
	} else {
		ws, ok := ttlConf.Workspaces[name]
		if !ok {
			return "", "", 0, fmt.Errorf("工作空间不存在: %s", name)
		}
		dbPath = ws.DbPath
		storageType = ws.StorageType
		if storageType == "" {
			storageType = DefaultStorageType
		}
	}

	var count int
	if dbPath != "" {
		if _, err := os.Stat(dbPath); err == nil {
			count = 1
		}
	}

	return dbPath, storageType, count, nil
}
