package cache

import (
	"backend/internal/store"
	"context"

	"github.com/stretchr/testify/mock"
)

func NewMockCache() Storage {
	return Storage{
		Users: &MockUserCache{},
	}
}

type MockUserCache struct {
	mock.Mock
}

func (s *MockUserCache) Get(ctx context.Context, userID int64) (*store.User, error) {
	args := s.Called(userID)
	return nil, args.Error(1)
}
func (s *MockUserCache) Set(ctx context.Context, user *store.User) error {
	args := s.Called(user)
	return args.Error(0)
}
