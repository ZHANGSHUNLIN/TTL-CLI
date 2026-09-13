package config

// TtlIni is the on-disk client configuration shape.
type TtlIni struct {
	StorageType string                     `ini:"storage_type"`
	DbPath      string                     `ini:"db_path"`
	Workspace   string                     `ini:"workspace"`
	BoltDB      BoltDBConfig               `ini:"bbolt"`
	Workspaces  map[string]WorkspaceConfig `ini:"-"`
}

type BoltDBConfig struct {
	Timeout int `ini:"timeout"`
}

type WorkspaceConfig struct {
	DbPath      string `ini:"db_path"`
	StorageType string `ini:"storage_type"`
}

type WorkspacesSection struct {
	Workspaces map[string]WorkspaceConfig `ini:"-"`
}
