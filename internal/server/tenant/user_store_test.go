package tenant

import (
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
	manager := NewStorageManager(t.TempDir())
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
}
