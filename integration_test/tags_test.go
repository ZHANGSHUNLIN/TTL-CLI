package integration_test

import (
	"os"
	"testing"

	"ttl-cli/internal/core/resource"
	storagebbolt "ttl-cli/internal/storage/bbolt"
)

var (
	testKey1 = resource.ValJsonKey{Key: "test-resource-1", Type: resource.ORIGIN}
	testVal1 = resource.ValJson{Val: "value1", Tag: []string{"work", "dev"}}
	testKey2 = resource.ValJsonKey{Key: "test-resource-2", Type: resource.ORIGIN}
	testVal2 = resource.ValJson{Val: "value2", Tag: []string{"work", "ci"}}
	testKey3 = resource.ValJsonKey{Key: "test-resource-3", Type: resource.ORIGIN}
	testVal3 = resource.ValJson{Val: "value3", Tag: []string{"deploy"}}
)

func TestTagsList(t *testing.T) {
	cleanup := setupTempStorage(t)
	defer cleanup()

	_ = testDB.SaveResource(testKey1, testVal1)
	_ = testDB.SaveResource(testKey2, testVal2)
	_ = testDB.SaveResource(testKey3, testVal3)

	stats, err := testDB.GetTagStats()
	if err != nil {
		t.Fatalf("GetTagStats failed: %v", err)
	}

	if len(stats) != 4 {
		t.Errorf("Expected 4 tags (work, dev, ci, deploy), got %d", len(stats))
	}

	tagMap := make(map[string]int)
	for _, stat := range stats {
		tagMap[stat.Tag] = stat.Count
	}

	if tagMap["work"] != 2 {
		t.Errorf("Expected work tag count 2, got %d", tagMap["work"])
	}
	if tagMap["dev"] != 1 {
		t.Errorf("Expected dev tag count 1, got %d", tagMap["dev"])
	}
	if tagMap["ci"] != 1 {
		t.Errorf("Expected ci tag count 1, got %d", tagMap["ci"])
	}
	if tagMap["deploy"] != 1 {
		t.Errorf("Expected deploy tag count 1, got %d", tagMap["deploy"])
	}
}

func TestTagsEmpty(t *testing.T) {
	cleanup := setupTempStorage(t)
	defer cleanup()

	stats, err := testDB.GetTagStats()
	if err != nil {
		t.Fatalf("GetTagStats failed: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(stats))
	}
}

func TestTagsSortOrder(t *testing.T) {
	cleanup := setupTempStorage(t)
	defer cleanup()

	_ = testDB.SaveResource(resource.ValJsonKey{Key: "z-test", Type: resource.ORIGIN}, resource.ValJson{Val: "value", Tag: []string{"zebra"}})
	_ = testDB.SaveResource(resource.ValJsonKey{Key: "a-test", Type: resource.ORIGIN}, resource.ValJson{Val: "value", Tag: []string{"alpha"}})

	stats, err := testDB.GetTagStats()
	if err != nil {
		t.Fatalf("GetTagStats failed: %v", err)
	}

	if len(stats) != 2 {
		t.Fatalf("Expected 2 tags, got %d", len(stats))
	}

	if stats[0].Tag != "alpha" {
		t.Errorf("First tag should be 'alpha', got '%s'", stats[0].Tag)
	}
	if stats[1].Tag != "zebra" {
		t.Errorf("Second tag should be 'zebra', got '%s'", stats[1].Tag)
	}
}

func TestTagsWithMultipleResources(t *testing.T) {
	cleanup := setupTempStorage(t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		key := resource.ValJsonKey{Key: "res-" + string(rune('a'+i)), Type: resource.ORIGIN}
		val := resource.ValJson{Val: "value", Tag: []string{"common"}}
		_ = testDB.SaveResource(key, val)
	}

	stats, err := testDB.GetTagStats()
	if err != nil {
		t.Fatalf("GetTagStats failed: %v", err)
	}

	if len(stats) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(stats))
	}

	if stats[0].Tag != "common" {
		t.Errorf("Expected tag 'common', got '%s'", stats[0].Tag)
	}

	if stats[0].Count != 5 {
		t.Errorf("Expected 5 resources, got %d", stats[0].Count)
	}
}

func TestTagsWithNoTagResources(t *testing.T) {
	cleanup := setupTempStorage(t)
	defer cleanup()

	_ = testDB.SaveResource(resource.ValJsonKey{Key: "no-tag-1", Type: resource.ORIGIN}, resource.ValJson{Val: "value", Tag: []string{}})
	_ = testDB.SaveResource(resource.ValJsonKey{Key: "no-tag-2", Type: resource.ORIGIN}, resource.ValJson{Val: "value", Tag: []string{}})

	stats, err := testDB.GetTagStats()
	if err != nil {
		t.Fatalf("GetTagStats failed: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(stats))
	}
}

func TestTagTypesExcluded(t *testing.T) {
	cleanup := setupTempStorage(t)
	defer cleanup()

	_ = testDB.SaveResource(testKey1, testVal1)
	_ = testDB.SaveResource(resource.ValJsonKey{Key: "tag-type", Type: resource.TAG, OriginKey: testKey1.Key}, resource.ValJson{Val: "value"})

	stats, err := testDB.GetTagStats()
	if err != nil {
		t.Fatalf("GetTagStats failed: %v", err)
	}

	if len(stats) != 2 {
		t.Errorf("Expected 2 tags (TAG type resources should be excluded, testVal1 has 'work' and 'dev'), got %d", len(stats))
	}
}

func testKey(prefix string, i int) resource.ValJsonKey {
	return resource.ValJsonKey{
		Key:  prefix + "-" + string(rune('a'+i)),
		Type: resource.ORIGIN,
	}
}

func TestGetTagStats_FileStorageOnly(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ttl-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testConf := tmpDir + "/test.conf"
	testDB := tmpDir + "/test.bbolt"

	confContent := "db_path = " + testDB + "\nstorage_type = bbolt\n"
	if err := os.WriteFile(testConf, []byte(confContent), 0644); err != nil {
		t.Fatal(err)
	}

	ls := storagebbolt.NewLocalStorage()
	ls.SetDBPath(testDB)
	if err := ls.Init(); err != nil {
		t.Fatal(err)
	}
	defer ls.Close()

	key1 := resource.ValJsonKey{Key: "file-test-1", Type: resource.ORIGIN}
	val1 := resource.ValJson{Val: "value1", Tag: []string{"file-tag"}}
	if err := ls.SaveResource(key1, val1); err != nil {
		t.Fatal(err)
	}

	stats, err := ls.GetTagStats()
	if err != nil {
		t.Fatalf("GetTagStats failed: %v", err)
	}

	if len(stats) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(stats))
	}

	if stats[0].Tag != "file-tag" {
		t.Errorf("Expected tag 'file-tag', got '%s'", stats[0].Tag)
	}

	if stats[0].Count != 1 {
		t.Errorf("Expected 1 resource, got %d", stats[0].Count)
	}
}
