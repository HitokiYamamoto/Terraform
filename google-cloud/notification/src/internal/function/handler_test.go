package function

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/config"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/repository"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- Mocks ---

// MockRepository は repository.Client の動作を模倣するモック
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetState(ctx context.Context, id string) (*repository.State, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.State), args.Error(1)
}

func (m *MockRepository) SaveState(ctx context.Context, id string, state *repository.State) error {
	args := m.Called(ctx, id, state)
	return args.Error(0)
}

// --- Tests ---

func TestBudgetAlertHandler_HandleBudgetAlert(t *testing.T) {
	// 共通の設定
	ctx := context.Background()
	cfg := &config.Config{ChannelName: "test-channel"}

	// 生存確認(Heartbeat)で誤爆しないように、現在時刻を用意
	now := time.Now()

	// テストケース定義
	tests := []struct {
		name string
		// Given
		lastState *repository.State // DBにある前回の状態
		alertData map[string]interface{}
		// When & Then
		expectNotify bool // Slack通知が飛ぶべきか
		expectSave   bool // DB保存が走るべきか
		expectError  bool // エラーが返るべきか
	}{
		{
			name: "【通知あり】前回50% -> 今回80% (上昇)",
			lastState: &repository.State{
				LastThreshold: 0.5,
				CurrentMonth:  "2026-02-01T00:00:00Z",
				LastHeartbeat: now,
			},
			alertData: map[string]interface{}{
				"budgetDisplayName":      "テスト予算",
				"alertThresholdExceeded": 0.8,
				"costAmount":             800.0,
				"budgetAmount":           1000.0,
				"currencyCode":           "USD",
				"costIntervalStart":      "2026-02-01T00:00:00Z",
			},
			expectNotify: true,
			expectSave:   true,
			expectError:  false,
		},
		{
			name: "【通知なし】前回50% -> 今回50% (変化なし)",
			lastState: &repository.State{
				LastThreshold: 0.5,
				CurrentMonth:  "2026-02-01T00:00:00Z",
				LastHeartbeat: now,
			},
			alertData: map[string]interface{}{
				"budgetDisplayName":      "テスト予算",
				"alertThresholdExceeded": 0.5,
				"costAmount":             500.0,
				"budgetAmount":           1000.0,
				"currencyCode":           "USD",
				"costIntervalStart":      "2026-02-01T00:00:00Z",
			},
			expectNotify: false,
			expectSave:   false,
			expectError:  false,
		},
		{
			name: "【通知あり】月が変わった場合 (リセット)",
			lastState: &repository.State{
				LastThreshold: 1.0,
				CurrentMonth:  "2026-01-01T00:00:00Z",
				LastHeartbeat: now,
			},
			alertData: map[string]interface{}{
				"budgetDisplayName":      "テスト予算",
				"alertThresholdExceeded": 0.1,
				"costAmount":             100.0,
				"budgetAmount":           1000.0,
				"currencyCode":           "USD",
				"costIntervalStart":      "2026-02-01T00:00:00Z",
			},
			expectNotify: true,
			expectSave:   true,
			expectError:  false,
		},
		{
			name:      "【通知あり】初回起動 (データなし)",
			lastState: nil,
			alertData: map[string]interface{}{
				"budgetDisplayName":      "テスト予算",
				"alertThresholdExceeded": 0.5,
				"costAmount":             500.0,
				"budgetAmount":           1000.0,
				"currencyCode":           "USD",
				"costIntervalStart":      "2026-02-01T00:00:00Z",
			},
			expectNotify: true,
			expectSave:   true,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Mocks
			mockSlack := slack.NewMockClient()
			mockRepo := new(MockRepository)

			// GetStateの挙動設定
			if tt.lastState == nil {
				// mock.Anything（第2引数のIDはなんでも良いという意味）
				mockRepo.On("GetState", ctx, mock.Anything).Return(nil, status.Error(codes.NotFound, "not found"))
			} else {
				mockRepo.On("GetState", ctx, mock.Anything).Return(tt.lastState, nil)
			}

			// SaveStateの挙動設定
			if tt.expectSave {
				mockRepo.On("SaveState", ctx, mock.Anything, mock.Anything).Return(nil)
			}

			// JSON作成
			jsonData, _ := json.Marshal(tt.alertData)
			message := PubSubMessage{Data: jsonData}

			// Handler作成
			handler := NewBudgetAlertHandler(mockSlack, mockRepo, cfg)

			// Execute
			err := handler.HandleBudgetAlert(ctx, message)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectNotify {
				assert.Equal(t, 1, len(mockSlack.Messages), "Slack通知件数不一致")
			} else {
				assert.Equal(t, 0, len(mockSlack.Messages), "Slack通知が送信されてはいけません")
			}

			if tt.expectSave {
				mockRepo.AssertNumberOfCalls(t, "SaveState", 1)
			} else {
				mockRepo.AssertNotCalled(t, "SaveState", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestNewBudgetAlertHandler(t *testing.T) {
	t.Run("ハンドラーが正しく作成される", func(t *testing.T) {
		mockSlack := slack.NewMockClient()
		mockRepo := new(MockRepository)
		cfg := &config.Config{ChannelName: "test-channel"}

		handler := NewBudgetAlertHandler(mockSlack, mockRepo, cfg)

		assert.NotNil(t, handler)
		assert.Equal(t, mockSlack, handler.slackClient)
		assert.Equal(t, cfg, handler.cfg)
	})
}
