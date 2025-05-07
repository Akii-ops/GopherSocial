package main

import (
	"backend/internal/store/cache"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	withRedis := config{
		redis: redisConfig{
			enabled: true,
		},
	}
	app := newTestApplication(t, withRedis)
	mux := app.mount()
	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {

		// check 401

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		mockCacheStore := app.cachestore.Users.(*cache.MockUserCache)

		mockCacheStore.On("Get", int64(1)).Return(nil, nil).Twice()
		mockCacheStore.On("Set", mock.Anything).Return(nil)
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.Calls = nil // Reset mock expectations
	})

	t.Run("should hit the cache first and if not exists it sets the user on the cache", func(t *testing.T) {
		mockCache := app.cachestore.Users.(*cache.MockUserCache)

		mockCache.On("Get", int64(1)).Return(nil, nil)
		mockCache.On("Get", int64(42)).Return(nil, nil)
		mockCache.On("Set", mock.Anything, mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCache.AssertNumberOfCalls(t, "Get", 2)

	})

	t.Run("should not hit the cache if cache is not enabled", func(t *testing.T) {
		withRedis := config{
			redis: redisConfig{
				enabled: false,
			},
		}

		app := newTestApplication(t, withRedis)
		mux := app.mount()

		mockCache := app.cachestore.Users.(*cache.MockUserCache)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCache.AssertNotCalled(t, "Get")
		mockCache.Calls = nil
	})
}
