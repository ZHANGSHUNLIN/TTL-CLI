package app

import (
	"errors"
	"reflect"
	"testing"

	"ttl-cli/internal/core/resource"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"
)

type resourceStorage struct {
	corestorage.Storage
	resources map[models.ValJsonKey]models.ValJson
	readErr   error
	writeErr  error
	updated   bool
}

func (s *resourceStorage) UpdateResource(key resource.Key, value resource.Value) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	s.updated = true
	value.CreatedAt = 30
	value.UpdatedAt = 40
	s.resources[key] = value
	return nil
}

func (s *resourceStorage) GetAllResources() (map[models.ValJsonKey]models.ValJson, error) {
	if s.readErr != nil {
		return nil, s.readErr
	}
	result := make(map[models.ValJsonKey]models.ValJson, len(s.resources))
	for key, value := range s.resources {
		result[key] = value
	}
	return result, nil
}

func (s *resourceStorage) SaveResource(key resource.Key, value resource.Value) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	if s.resources == nil {
		s.resources = map[models.ValJsonKey]models.ValJson{}
	}
	value.CreatedAt = 10
	value.UpdatedAt = 20
	s.resources[key] = value
	return nil
}

func (s *resourceStorage) DeleteResource(key resource.Key) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	delete(s.resources, key)
	return nil
}

func TestService_FindResourcesUsesKeyAndTagMatching(t *testing.T) {
	storage := &resourceStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "alpha", Type: models.ORIGIN}:                      {Val: "one", Tag: []string{"work"}, CreatedAt: 1},
		{Key: "beta", Type: models.ORIGIN}:                       {Val: "alpha in value", Tag: []string{"other"}, CreatedAt: 3},
		{Key: "gamma", Type: models.ORIGIN}:                      {Val: "three", Tag: []string{"alpha-tag"}, CreatedAt: 2},
		{Key: "alpha-tag", Type: models.TAG, OriginKey: "gamma"}: {Val: "gamma", CreatedAt: 4},
	}}
	service := NewService(storage)

	matches, err := service.FindResources("alpha")
	if err != nil {
		t.Fatalf("FindResources() error = %v", err)
	}
	got := []string{matches[0].Key.Key, matches[1].Key.Key, matches[2].Key.Key}
	want := []string{"beta", "gamma", "alpha"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindResources() keys = %v, want %v", got, want)
	}
}

func TestService_FindResourcesWithOptionsDefaultsToKeyAndCanIncludeValue(t *testing.T) {
	service := NewService(&resourceStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "alpha-key", Type: models.ORIGIN}: {Val: "first value"},
		{Key: "beta-key", Type: models.ORIGIN}:  {Val: "alpha in value"},
	}})

	keyMatches, err := service.FindResourcesWithOptions("alpha", SearchOptions{})
	if err != nil || len(keyMatches) != 1 || keyMatches[0].Key.Key != "alpha-key" {
		t.Fatalf("key-only matches = %+v, err = %v", keyMatches, err)
	}
	valueMatches, err := service.FindResourcesWithOptions("alpha", SearchOptions{IncludeValue: true})
	if err != nil || len(valueMatches) != 2 {
		t.Fatalf("key/value matches = %+v, err = %v", valueMatches, err)
	}
}

func TestService_FindResourcesWithOptionsCanIncludeTagsWithoutValue(t *testing.T) {
	service := NewService(&resourceStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "deployment-note", Type: models.ORIGIN}: {Val: "secret", Tag: []string{"production"}},
	}})

	matches, err := service.FindResourcesWithOptions("production", SearchOptions{IncludeTags: true})
	if err != nil || len(matches) != 1 || matches[0].Key.Key != "deployment-note" {
		t.Fatalf("tag matches = %+v, err = %v", matches, err)
	}
}

func TestService_FindResourcesEmptyQueryListsOriginResources(t *testing.T) {
	service := NewService(&resourceStorage{resources: map[models.ValJsonKey]models.ValJson{
		{Key: "note", Type: models.ORIGIN}:                {Val: "value"},
		{Key: "tag", Type: models.TAG, OriginKey: "note"}: {Val: "note"},
	}})

	matches, err := service.FindResources("")
	if err != nil || len(matches) != 1 || matches[0].Key.Key != "note" {
		t.Fatalf("FindResources(empty) = %+v, %v", matches, err)
	}
}

func TestService_FindResourcesReturnsTypedNotFound(t *testing.T) {
	service := NewService(&resourceStorage{resources: map[models.ValJsonKey]models.ValJson{}})

	_, err := service.FindResources("missing")
	if kind, ok := ErrorKindOf(err); !ok || kind != ErrorNotFound {
		t.Fatalf("FindResources() error = %v, kind = %q, want %q", err, kind, ErrorNotFound)
	}
}

func TestService_CreateResourceRejectsDuplicate(t *testing.T) {
	key := models.ValJsonKey{Key: "duplicate", Type: models.ORIGIN}
	service := NewService(&resourceStorage{resources: map[models.ValJsonKey]models.ValJson{key: {Val: "old"}}})

	_, err := service.CreateResource("duplicate", "new", nil)
	if kind, ok := ErrorKindOf(err); !ok || kind != ErrorConflict {
		t.Fatalf("CreateResource() error = %v, kind = %q, want %q", err, kind, ErrorConflict)
	}
}

func TestService_UpdateResourceValuePreservesTags(t *testing.T) {
	key := models.ValJsonKey{Key: "note", Type: models.ORIGIN}
	storage := &resourceStorage{resources: map[models.ValJsonKey]models.ValJson{key: {Val: "old", Tag: []string{"work"}}}}
	service := NewService(storage)

	updated, err := service.UpdateResourceValue("note", "new")
	if err != nil {
		t.Fatalf("UpdateResourceValue() error = %v", err)
	}
	if updated.Value.Val != "new" || !reflect.DeepEqual(updated.Value.Tag, []string{"work"}) {
		t.Fatalf("UpdateResourceValue() = %+v", updated.Value)
	}
	if !storage.updated || updated.Value.CreatedAt != 30 || updated.Value.UpdatedAt != 40 {
		t.Fatalf("UpdateResourceValue() did not use UpdateResource: updated=%v value=%+v", storage.updated, updated.Value)
	}
}

func TestService_CreateResourceWrapsStorageFailure(t *testing.T) {
	writeErr := errors.New("disk full")
	service := NewService(&resourceStorage{resources: map[models.ValJsonKey]models.ValJson{}, writeErr: writeErr})

	_, err := service.CreateResource("note", "value", nil)
	if kind, ok := ErrorKindOf(err); !ok || kind != ErrorSystem || !errors.Is(err, writeErr) {
		t.Fatalf("CreateResource() error = %v, kind = %q", err, kind)
	}
}

func TestService_DeleteResourceReportsCleanupFailuresAndDeletes(t *testing.T) {
	key := models.ValJsonKey{Key: "note", Type: models.ORIGIN}
	cleanupErr := errors.New("cleanup failed")
	storage := &resourceStorage{resources: map[models.ValJsonKey]models.ValJson{key: {Val: "value"}}}
	service := NewService(&deleteStorage{resourceStorage: storage, cleanupErr: cleanupErr})

	result, err := service.DeleteResourceWithCleanup("note")
	if err != nil {
		t.Fatalf("DeleteResourceWithCleanup() error = %v", err)
	}
	if !errors.Is(result.HistoryCleanupError, cleanupErr) || !errors.Is(result.AuditCleanupError, cleanupErr) {
		t.Fatalf("DeleteResourceWithCleanup() cleanup = %+v", result)
	}
	if _, exists := storage.resources[key]; exists {
		t.Fatal("resource was not deleted")
	}
}

type deleteStorage struct {
	*resourceStorage
	cleanupErr error
}

func (s *deleteStorage) DeleteHistoryRecords(string) error { return s.cleanupErr }
func (s *deleteStorage) DeleteAuditRecords(string) error   { return s.cleanupErr }
