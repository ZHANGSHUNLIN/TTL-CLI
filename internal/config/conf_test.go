package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfFile_RemoteProfiles(t *testing.T) {
	confFile := filepath.Join(t.TempDir(), "ttl.ini")
	content := `[storage]
type = local
remote = primary

[remotes.primary]
url = http://primary.example
account = alice
credential_env = TTL_PRIMARY_KEY

[remotes.backup]
url = http://backup.example
credential_env = TTL_BACKUP_KEY

[workspaces.demo]
storage_type = cloud
remote = backup
`
	if err := os.WriteFile(confFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	conf, err := GetTtlConfFromFile(confFile)
	if err != nil {
		t.Fatal(err)
	}
	if conf.StorageType != "local" || len(conf.Remotes) != 2 {
		t.Fatalf("config = %+v, want local and two remotes", conf)
	}
	if conf.Workspaces["demo"].RemoteName != "backup" {
		t.Fatalf("workspace remote = %q, want backup", conf.Workspaces["demo"].RemoteName)
	}
	name, remote, err := LoadRemote(confFile)
	if err != nil || name != "primary" || remote.URL != "http://primary.example" {
		t.Fatalf("LoadRemote() = %q, %+v, %v", name, remote, err)
	}
}

func TestLoadRemote_UsesActiveWorkspace(t *testing.T) {
	confFile := filepath.Join(t.TempDir(), "ttl.ini")
	content := `workspace = demo

[storage]
type = cloud
remote = primary

[remotes.primary]
url = http://primary.example

[remotes.backup]
url = http://backup.example

[workspaces.demo]
storage_type = cloud
remote = backup
`
	if err := os.WriteFile(confFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	name, remote, err := LoadRemote(confFile)
	if err != nil {
		t.Fatal(err)
	}
	if name != "backup" || remote.URL != "http://backup.example" {
		t.Fatalf("LoadRemote() = %q, %+v", name, remote)
	}
}

func TestValidateStorageType_RejectsRemovedModes(t *testing.T) {
	for _, value := range []string{"sqlite", "bbolt", "sync"} {
		if err := ValidateStorageType(value); err == nil {
			t.Errorf("ValidateStorageType(%q) returned nil", value)
		}
	}
	for _, value := range []string{"local", "cloud", ""} {
		if err := ValidateStorageType(value); err != nil {
			t.Errorf("ValidateStorageType(%q) = %v", value, err)
		}
	}
}

func TestGetTtlConf(t *testing.T) {
	tests := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "测试默认配置",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, err := GetTtlConf()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("GetTtlConf() error = %v, want nil", err)
				}
				// 验证配置是否包含必要的字段
				// DbPath 是动态生成的，首次创建时为空是正常的
				if conf.StorageType == "" {
					t.Errorf("GetTtlConf() returned config with empty StorageType")
				}
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "测试配置初始化",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTtlConf()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Failed to get config: %v", err)
				}
			}
		})
	}
}
