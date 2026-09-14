package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_Ready(t *testing.T) {
	handler := HealthHandler(func() bool { return true })
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Body.String(); got != "{\"status\":\"ok\"}\n" {
		t.Fatalf("body = %q, want ready response", got)
	}
}

func TestHealthHandler_NotReady(t *testing.T) {
	handler := HealthHandler(func() bool { return false })
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusServiceUnavailable)
	}
	if got := res.Body.String(); got != "{\"status\":\"starting\"}\n" {
		t.Fatalf("body = %q, want starting response", got)
	}
}

func TestHealthHandler_OnlyAllowsGet(t *testing.T) {
	handler := HealthHandler(func() bool { return true })
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
	if got := res.Header().Get("Allow"); got != http.MethodGet {
		t.Fatalf("Allow = %q, want GET", got)
	}
}
