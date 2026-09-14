package config

// TtlIni is the on-disk client configuration shape.
type TtlIni struct {
	StorageType string                     `ini:"storage_type"`
	DbPath      string                     `ini:"db_path"`
	Workspace   string                     `ini:"workspace"`
	RemoteName  string                     `ini:"remote"`
	BoltDB      BoltDBConfig               `ini:"bbolt"`
	Workspaces  map[string]WorkspaceConfig `ini:"-"`
	Remotes     map[string]RemoteConfig    `ini:"-"`
}

type BoltDBConfig struct {
	Timeout int `ini:"timeout"`
}

type WorkspaceConfig struct {
	DbPath      string `ini:"db_path"`
	StorageType string `ini:"storage_type"`
	RemoteName  string `ini:"remote"`
}

// RemoteConfig identifies a configured cloud endpoint without storing a secret.
type RemoteConfig struct {
	URL           string `ini:"url"`
	Account       string `ini:"account"`
	CredentialEnv string `ini:"credential_env"`
}

type WorkspacesSection struct {
	Workspaces map[string]WorkspaceConfig `ini:"-"`
}
