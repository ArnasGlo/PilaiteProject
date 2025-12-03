package mocks

import "context"

type MockSessionManager struct {
	GetIntFunc    func(ctx context.Context, key string) int
	GetStringFunc func(ctx context.Context, key string) string
}

func (m *MockSessionManager) GetInt(ctx context.Context, key string) int {
	return m.GetIntFunc(ctx, key)
}

func (m *MockSessionManager) GetString(ctx context.Context, key string) string {
	return m.GetStringFunc(ctx, key)
}
