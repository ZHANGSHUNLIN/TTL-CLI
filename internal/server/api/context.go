package api

import (
	"context"
	"net/http"
	corestorage "ttl-cli/internal/core/storage"
)

type contextKey string

const (
	contextKeyUserID  contextKey = "user_id"
	contextKeyStorage contextKey = "storage"
)

func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value(contextKeyUserID).(string)
	return userID
}

func GetStorage(r *http.Request) corestorage.Storage {
	storage, _ := r.Context().Value(contextKeyStorage).(corestorage.Storage)
	return storage
}

func withUserStorage(ctx context.Context, userID string, storage corestorage.Storage) context.Context {
	ctx = context.WithValue(ctx, contextKeyUserID, userID)
	return context.WithValue(ctx, contextKeyStorage, storage)
}
