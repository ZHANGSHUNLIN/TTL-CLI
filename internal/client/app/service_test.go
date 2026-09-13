package app

import (
	"context"
	"errors"
	"testing"

	"ttl-cli/internal/core/resource"
	corestorage "ttl-cli/internal/core/storage"
)

type trackingStorage struct {
	corestorage.Storage
	closed     bool
	historyErr error
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
