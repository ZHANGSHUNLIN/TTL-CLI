package integration_test

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"slices"
	"testing"
	"time"

	corestorage "ttl-cli/internal/core/storage"
	storagebbolt "ttl-cli/internal/storage/bbolt"
	storagesqlite "ttl-cli/internal/storage/sqlite"
	"ttl-cli/models"

	"go.etcd.io/bbolt"
)

const legacyTimestamp int64 = 1700000000

type legacyResourceKey struct {
	Key       string `json:"key"`
	Type      int    `json:"type"`
	OriginKey string `json:"originKey"`
}

func TestLocalStorage_UpdatePreservesResourceMetadata(t *testing.T) {
	for _, test := range []struct {
		name string
		open func(*testing.T) corestorage.Storage
	}{
		{name: "bbolt", open: openTestBbolt},
		{name: "sqlite", open: openTestSQLite},
	} {
		t.Run(test.name, func(t *testing.T) {
			storage := test.open(t)
			key := models.ValJsonKey{Key: "note", Type: models.ORIGIN}
			if err := storage.SaveResource(key, models.ValJson{Val: "before", Tag: []string{"work"}}); err != nil {
				t.Fatalf("SaveResource() error = %v", err)
			}
			before, err := storage.GetAllResources()
			if err != nil {
				t.Fatalf("GetAllResources(before) error = %v", err)
			}
			time.Sleep(1100 * time.Millisecond)
			value := before[key]
			value.Val = "after"
			if err := storage.UpdateResource(key, value); err != nil {
				t.Fatalf("UpdateResource() error = %v", err)
			}
			after, err := storage.GetAllResources()
			if err != nil {
				t.Fatalf("GetAllResources(after) error = %v", err)
			}
			got := after[key]
			if got.CreatedAt != value.CreatedAt || got.UpdatedAt <= value.UpdatedAt || got.Val != "after" || !slices.Equal(got.Tag, []string{"work"}) {
				t.Fatalf("updated metadata = %+v, before = %+v", got, value)
			}
		})
	}
}

func openTestBbolt(t *testing.T) corestorage.Storage {
	t.Helper()
	storage := storagebbolt.NewLocalStorage()
	storage.SetDBPath(filepath.Join(t.TempDir(), "data.bbolt"))
	if err := storage.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	t.Cleanup(func() { _ = storage.Close() })
	return storage
}

func openTestSQLite(t *testing.T) corestorage.Storage {
	t.Helper()
	storage := storagesqlite.NewSQLiteStorage()
	storage.SetDBPath(filepath.Join(t.TempDir(), "data.db"))
	if err := storage.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	t.Cleanup(func() { _ = storage.Close() })
	return storage
}

type legacyResourceValue struct {
	Val       string   `json:"val"`
	Tag       []string `json:"tag"`
	CreatedAt int64    `json:"createdAt"`
	UpdatedAt int64    `json:"updatedAt"`
}

type legacyAuditRecord struct {
	ResourceKey string `json:"resourceKey"`
	Operation   string `json:"operation"`
	Timestamp   int64  `json:"timestamp"`
	Count       int    `json:"count"`
}

type legacyHistoryRecord struct {
	ID          int64  `json:"id"`
	ResourceKey string `json:"resourceKey"`
	Operation   string `json:"operation"`
	Timestamp   int64  `json:"timestamp"`
	TimeStr     string `json:"timeStr"`
	Command     string `json:"command"`
	Args        string `json:"args"`
}

type legacyLogRecord struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
	Date      string   `json:"date"`
}

func TestBboltStorage_LegacyDataCompatibility(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "legacy.bbolt")
	createLegacyBboltFixture(t, dbPath)

	storage := storagebbolt.NewLocalStorage()
	storage.SetDBPath(dbPath)
	if err := storage.Init(); err != nil {
		t.Fatalf("open legacy bbolt data: %v", err)
	}
	t.Cleanup(func() { _ = storage.Close() })

	assertLegacyStorageData(t, storage)
}

func TestSQLiteStorage_LegacyDataCompatibility(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "legacy.db")
	createLegacySQLiteFixture(t, dbPath)

	storage := storagesqlite.NewSQLiteStorage()
	storage.SetDBPath(dbPath)
	if err := storage.Init(); err != nil {
		t.Fatalf("open legacy sqlite data: %v", err)
	}
	t.Cleanup(func() { _ = storage.Close() })

	assertLegacyStorageData(t, storage)
}

func createLegacyBboltFixture(t *testing.T, path string) {
	t.Helper()

	database, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		t.Fatalf("create legacy bbolt fixture: %v", err)
	}

	resourceKey := mustJSON(t, legacyResourceKey{Key: "legacy-resource", Type: 0})
	fixtures := []struct {
		bucket string
		key    []byte
		value  []byte
	}{
		{"resources", resourceKey, mustJSON(t, legacyResourceValue{Val: "legacy-value", Tag: []string{"legacy", "fixture"}, CreatedAt: legacyTimestamp, UpdatedAt: legacyTimestamp})},
		{"resources", mustJSON(t, legacyResourceKey{Key: "legacy-tag", Type: 1, OriginKey: "legacy-resource"}), mustJSON(t, legacyResourceValue{Val: "legacy-resource", CreatedAt: legacyTimestamp, UpdatedAt: legacyTimestamp})},
		{"audit", []byte("legacy-resource_add_1700000000"), mustJSON(t, legacyAuditRecord{ResourceKey: "legacy-resource", Operation: "add", Timestamp: legacyTimestamp, Count: 2})},
		{"history", []byte("1700000000_7"), mustJSON(t, legacyHistoryRecord{ID: 7, ResourceKey: "legacy-resource", Operation: "add", Timestamp: legacyTimestamp, TimeStr: "2023-11-14 22:13:20", Command: "ttl add"})},
		{"logs", []byte("2023-11-14_9"), mustJSON(t, legacyLogRecord{ID: 9, Content: "legacy log", Tags: []string{"legacy"}, CreatedAt: "2023-11-14 22:13:20", Date: "2023-11-14"})},
	}

	if err := database.Update(func(tx *bbolt.Tx) error {
		for _, fixture := range fixtures {
			bucket, err := tx.CreateBucketIfNotExists([]byte(fixture.bucket))
			if err != nil {
				return err
			}
			if err := bucket.Put(fixture.key, fixture.value); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		_ = database.Close()
		t.Fatalf("write legacy bbolt fixture: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close legacy bbolt fixture: %v", err)
	}
}

func createLegacySQLiteFixture(t *testing.T, path string) {
	t.Helper()

	database, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("create legacy sqlite fixture: %v", err)
	}

	statements := []string{
		`CREATE TABLE resources (key TEXT NOT NULL, type TEXT NOT NULL, origin_key TEXT, value TEXT NOT NULL, tags TEXT, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL, PRIMARY KEY (key, type))`,
		`CREATE TABLE audit (id INTEGER PRIMARY KEY AUTOINCREMENT, resource_key TEXT NOT NULL, operation TEXT NOT NULL, timestamp INTEGER NOT NULL, count INTEGER NOT NULL)`,
		`CREATE TABLE history (id INTEGER PRIMARY KEY, resource_key TEXT, operation TEXT NOT NULL, timestamp INTEGER NOT NULL, time_str TEXT NOT NULL, command TEXT NOT NULL)`,
		`CREATE TABLE logs (id INTEGER PRIMARY KEY, content TEXT NOT NULL, tags TEXT, created_at TEXT NOT NULL, date TEXT NOT NULL)`,
		`INSERT INTO resources (key, type, origin_key, value, tags, created_at, updated_at) VALUES ('legacy-resource', 'ORIGIN', '', 'legacy-value', '["legacy","fixture"]', 1700000000, 1700000000)`,
		`INSERT INTO resources (key, type, origin_key, value, tags, created_at, updated_at) VALUES ('legacy-tag', 'TAG', 'legacy-resource', 'legacy-resource', 'null', 1700000000, 1700000000)`,
		`INSERT INTO audit (resource_key, operation, timestamp, count) VALUES ('legacy-resource', 'add', 1700000000, 2)`,
		`INSERT INTO history (id, resource_key, operation, timestamp, time_str, command) VALUES (7, 'legacy-resource', 'add', 1700000000, '2023-11-14 22:13:20', 'ttl add')`,
		`INSERT INTO logs (id, content, tags, created_at, date) VALUES (9, 'legacy log', '["legacy"]', '2023-11-14 22:13:20', '2023-11-14')`,
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			_ = database.Close()
			t.Fatalf("write legacy sqlite fixture: %v", err)
		}
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close legacy sqlite fixture: %v", err)
	}
}

func assertLegacyStorageData(t *testing.T, storage corestorage.Storage) {
	t.Helper()

	key := models.ValJsonKey{Key: "legacy-resource", Type: models.ORIGIN}
	resources, err := storage.GetAllResources()
	if err != nil {
		t.Fatalf("read legacy resources: %v", err)
	}
	value, ok := resources[key]
	if !ok {
		t.Fatal("legacy resource was not readable")
	}
	if value.Val != "legacy-value" || !slices.Equal(value.Tag, []string{"legacy", "fixture"}) || value.CreatedAt != legacyTimestamp || value.UpdatedAt != legacyTimestamp {
		t.Fatalf("legacy resource changed: %+v", value)
	}
	tagKey := models.ValJsonKey{Key: "legacy-tag", Type: models.TAG, OriginKey: key.Key}
	if tag, ok := resources[tagKey]; !ok || tag.Val != key.Key {
		t.Fatalf("legacy tag resource changed: %+v", tag)
	}

	audit, err := storage.GetAllAuditRecords()
	if err != nil {
		t.Fatalf("read legacy audit: %v", err)
	}
	if len(audit) != 1 || audit[0].ResourceKey != key.Key || audit[0].Operation != "add" || audit[0].Count != 2 {
		t.Fatalf("legacy audit changed: %+v", audit)
	}

	history, err := storage.GetAllHistoryRecords()
	if err != nil {
		t.Fatalf("read legacy history: %v", err)
	}
	if len(history) != 1 || history[0].ID != 7 || history[0].ResourceKey != key.Key || history[0].Command != "ttl add" {
		t.Fatalf("legacy history changed: %+v", history)
	}

	logs, err := storage.GetLogRecords("", "")
	if err != nil {
		t.Fatalf("read legacy logs: %v", err)
	}
	if len(logs) != 1 || logs[0].ID != 9 || logs[0].Content != "legacy log" || !slices.Equal(logs[0].Tags, []string{"legacy"}) {
		t.Fatalf("legacy logs changed: %+v", logs)
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal legacy fixture: %v", err)
	}
	return data
}
