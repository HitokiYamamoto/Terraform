package secret

import (
	"context"
)

// MockManager はテスト用のモックSecret Manager
type MockManager struct {
	Secrets map[string]string
	Error   error
}

// NewMockManager は新しいモックSecret Managerを作成する
func NewMockManager() *MockManager {
	return &MockManager{
		Secrets: make(map[string]string),
	}
}

// GetSecret はモックからシークレットを取得する
func (m *MockManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	if secret, ok := m.Secrets[secretName]; ok {
		return secret, nil
	}

	return "", nil
}
