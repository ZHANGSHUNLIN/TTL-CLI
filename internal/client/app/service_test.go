package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"ttl-cli/internal/client/remote"
	"ttl-cli/internal/core/resource"
	corestorage "ttl-cli/internal/core/storage"
	storagesqlite "ttl-cli/internal/storage/sqlite"
)

type trackingStorage struct {
	corestorage.Storage
	closed     bool
	historyErr error
}

func TestOpenStorage_LocalUsesSQLite(t *testing.T) {
	confFile := filepath.Join(t.TempDir(), "ttl.ini")
	dbPath := filepath.Join(t.TempDir(), "data.sqlite")
	if err := os.WriteFile(confFile, []byte("[storage]\ntype = local\npath = "+dbPath+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	storage, err := OpenStorage("local", "", "", 0, confFile)
	if err != nil {
		t.Fatalf("OpenStorage(local): %v", err)
	}
	defer storage.Close()
	if _, ok := storage.(*storagesqlite.SQLiteStorage); !ok {
		t.Fatalf("storage type = %T, want SQLiteStorage", storage)
	}
}

func TestOpenStorage_UsesActiveRemoteProfile(t *testing.T) {
	confFile := filepath.Join(t.TempDir(), "ttl.ini")
	content := `[storage]
type = cloud
remote = primary

[remotes.primary]
url = http://primary.example
credential_env = TTL_TEST_REMOTE_KEY
`
	if err := os.WriteFile(confFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TTL_TEST_REMOTE_KEY", "test-key")
	storage, err := OpenStorage("cloud", "", "", 0, confFile)
	if err != nil {
		t.Fatalf("OpenStorage(cloud): %v", err)
	}
	defer storage.Close()
	if _, ok := storage.(*remote.Storage); !ok {
		t.Fatalf("storage type = %T, want remote.Storage", storage)
	}
}

func TestOpenStorage_RejectsRemovedModes(t *testing.T) {
	for _, mode := range []string{"sqlite", "bbolt", "sync"} {
		if _, err := OpenStorage(mode, "", "", 0, ""); err == nil {
			t.Errorf("OpenStorage(%q) returned nil error", mode)
		}
	}
}

func (s *trackingStorage) SaveHistoryRecord(resource.HistoryRecord) error {
	return s.historyErr
}

func (s *trackingStorage) Close() error {
	s.closed = true
	return nil
}

func TestServiceFromContext_MissingService(t *testing.T) {
	if _, err := ServiceFromContext(context.Background()); err == nil {
		t.Fatal("expected missing service error")
	}
}

func TestService_CloseDelegatesToStorage(t *testing.T) {
	storage := &trackingStorage{}
	service := NewService(storage)

	if err := service.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !storage.closed {
		t.Fatal("expected storage to be closed")
	}
}

func TestService_EmptyStorageReturnsInitializationError(t *testing.T) {
	service := NewService(nil)
	if _, err := service.GetAllResources(); err == nil {
		t.Fatal("expected storage not initialized error")
	}
}

func TestService_RecordCommandHistoryPropagatesStorageError(t *testing.T) {
	storageErr := errors.New("history write failed")
	service := NewService(&trackingStorage{historyErr: storageErr})

	if err := service.RecordCommandHistory("add", "key", false); !errors.Is(err, storageErr) {
		t.Fatalf("RecordCommandHistory() error = %v, want %v", err, storageErr)
	}
}
