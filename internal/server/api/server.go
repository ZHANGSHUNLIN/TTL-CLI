package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/internal/server/tenant"
)

func StartServer(port int, dataDir string) error {
	userStore := tenant.NewUserStore(filepath.Join(dataDir, "users.json"))
	if err := userStore.Load(); err != nil {
		return fmt.Errorf("failed to load user data: %w", err)
	}

	users := userStore.ListUsers()
	if len(users) == 0 {
		fmt.Println("Warning: No users, please run 'ttl server user add' to create users first")
	}

	tenantMgr := tenant.NewStorageManager(filepath.Join(dataDir, "tenants"))

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/resources", ResourcesHandler)
	mux.HandleFunc("/api/v1/resources/", ResourceHandler)
	mux.HandleFunc("/api/v1/audit/stats", AuditStatsHandler)
	mux.HandleFunc("/api/v1/history", HistoryHandler)

	handler := MultiTenantAuthMiddleware(userStore, tenantMgr, mux)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("ttl server started, listening on %s\n", addr)
	fmt.Printf("  REST API: http://localhost:%d/api/v1/\n", port)
	fmt.Printf("  Data dir: %s\n", dataDir)
	fmt.Printf("  User count: %d\n", len(users))

	return http.ListenAndServe(addr, handler)
}

// NewHandler builds the API router with an explicit storage dependency.
func NewHandler(storage corestorage.Storage) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/resources", ResourcesHandler)
	mux.HandleFunc("/api/v1/resources/", ResourceHandler)
	mux.HandleFunc("/api/v1/audit/stats", AuditStatsHandler)
	mux.HandleFunc("/api/v1/history", HistoryHandler)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := withUserStorage(r.Context(), "", storage)
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}
