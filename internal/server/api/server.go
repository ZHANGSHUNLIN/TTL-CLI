package api

import (
	"net/http"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/internal/server/tenant"
)

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

// NewServerHandler builds the production router with tenant authentication and
// an unauthenticated readiness endpoint.
func NewServerHandler(userStore *tenant.UserStore, tenantMgr *tenant.StorageManager, ready func() bool) http.Handler {
	protected := http.NewServeMux()
	protected.HandleFunc("/api/v1/resources", ResourcesHandler)
	protected.HandleFunc("/api/v1/resources/", ResourceHandler)
	protected.HandleFunc("/api/v1/audit/stats", AuditStatsHandler)
	protected.HandleFunc("/api/v1/history", HistoryHandler)

	root := http.NewServeMux()
	root.Handle("/healthz", HealthHandler(ready))
	root.Handle("/", MultiTenantAuthMiddleware(userStore, tenantMgr, protected))
	return root
}

// HealthHandler reports process readiness without exposing tenant data.
func HealthHandler(ready func() bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if ready != nil && ready() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}` + "\n"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"starting"}` + "\n"))
	})
}
