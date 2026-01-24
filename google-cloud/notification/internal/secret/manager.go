package secret

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Manager はSecret Managerからシークレットを取得するインターフェース
type Manager interface {
	GetSecret(ctx context.Context, secretName string) (string, error)
}

// GCPSecretManager はGCPのSecret Managerクライアント
type GCPSecretManager struct {
	projectID string
	client    *secretmanager.Client
}

// NewManager は新しいSecret Managerクライアントを作成する
func NewManager(ctx context.Context, projectID string) (*GCPSecretManager, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}

	return &GCPSecretManager{
		projectID: projectID,
		client:    client,
	}, nil
}

// GetSecret はSecret Managerからシークレットを取得する
func (m *GCPSecretManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", m.projectID, secretName)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := m.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret %s: %w", secretName, err)
	}

	return string(result.Payload.Data), nil
}

// Close はクライアントを閉じる
func (m *GCPSecretManager) Close() error {
	return m.client.Close()
}
