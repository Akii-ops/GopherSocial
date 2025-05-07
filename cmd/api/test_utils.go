package main

import (
	"backend/internal/auth"
	"backend/internal/store"
	"backend/internal/store/cache"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {

	t.Helper()
	logger := zap.Must(zap.NewProduction()).Sugar()
	// logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockCache()
	authenticator := &auth.TestAuthenticator{}
	return &application{
		config:        cfg,
		logger:        logger,
		store:         mockStore,
		cachestore:    mockCache,
		authenticator: authenticator,
	}
}

func executeRequest(req *http.Request, mux *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expect int, got int) {
	if expect != got {
		t.Errorf("expected the response code to be %v, but got %v", expect, got)

	}
}
