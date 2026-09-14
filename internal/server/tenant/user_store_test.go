package tenant

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUserStore_LifecycleAndPersistence(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "nested", "users.json")
	store := NewUserStore(filePath)
	if err := store.Load(); err != nil {
		t.Fatalf("Load empty store: %v", err)
	}
	user, err := store.AddUser("alice", "Alice")
	if err != nil {
		t.Fatalf("AddUser: %v", err)
	}
	if user.APIKey == "" || !user.Active {
		t.Fatalf("created user = %+v", user)
	}
	if err := store.SetActive("alice", false); err != nil {
		t.Fatalf("SetActive: %v", err)
	}
	if _, err := store.ResetKey("alice"); err != nil {
		t.Fatalf("ResetKey: %v", err)
	}

	reloaded := NewUserStore(filePath)
	if err := reloaded.Load(); err != nil {
		t.Fatalf("reload store: %v", err)
	}
	if found := reloaded.FindByID("alice"); found == nil || found.Active {
		t.Fatalf("reloaded user = %+v", found)
	}
	if err := reloaded.DeleteUser("alice"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
}

func TestStorageManager_IsolatesTenants(t *testing.T) {
	dataDir := t.TempDir()
	manager := NewStorageManager(dataDir)
	defer manager.CloseAll()

	alice, err := manager.GetStorage("alice")
	if err != nil {
		t.Fatalf("alice storage: %v", err)
	}
	bob, err := manager.GetStorage("bob")
	if err != nil {
		t.Fatalf("bob storage: %v", err)
	}
	if alice == bob {
		t.Fatal("different users must not share a storage instance")
	}
	if _, err := os.Stat(filepath.Join(dataDir, "alice", "data.sqlite")); err != nil {
		t.Fatalf("alice sqlite file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "bob", "data.sqlite")); err != nil {
		t.Fatalf("bob sqlite file: %v", err)
	}
}

func TestStorageManager_RemoveStorage_IsolatesData(t *testing.T) {
	dataDir := t.TempDir()
	manager := NewStorageManager(dataDir)
	if _, err := manager.GetStorage("alice"); err != nil {
		t.Fatal(err)
	}
	if err := manager.RemoveStorage("alice"); err != nil {
		t.Fatalf("RemoveStorage: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "alice")); !os.IsNotExist(err) {
		t.Fatalf("alice directory still present: %v", err)
	}
	entries, err := os.ReadDir(filepath.Join(dataDir, ".deleted"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("deleted entries = %d, err = %v", len(entries), err)
	}
	_ = manager.CloseAll()
}
