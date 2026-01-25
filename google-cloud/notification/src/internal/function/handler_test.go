package function

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/config"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/slack"
	"github.com/stretchr/testify/assert"
)

func TestBudgetAlertHandler_HandleBudgetAlert(t *testing.T) {
	t.Run("正しいメッセージの場合、Slackに通知が送信される", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		mockSlack := slack.NewMockClient()
		cfg := &config.Config{
			ChannelName: "test-channel",
		}
		handler := NewBudgetAlertHandler(mockSlack, cfg)

		// 予算アラートのデータを作成
		alertData := map[string]interface{}{
			"budgetDisplayName":      "テスト予算",
			"alertThresholdExceeded": 0.8,
			"costAmount":             800.0,
			"budgetAmount":           1000.0,
			"currencyCode":           "USD",
		}
		jsonData, _ := json.Marshal(alertData)

		message := PubSubMessage{
			Data: jsonData,
		}

		// Execute
		err := handler.HandleBudgetAlert(ctx, message)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, len(mockSlack.Messages))
		lastMsg := mockSlack.GetLastMessage()
		assert.NotNil(t, lastMsg)
		assert.Equal(t, "test-channel", lastMsg.Channel)
		assert.Contains(t, lastMsg.Text, "テスト予算")
		assert.Contains(t, lastMsg.Text, "80.00%")
	})

	t.Run("不正なメッセージの場合、エラーが返される", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		mockSlack := slack.NewMockClient()
		cfg := &config.Config{
			ChannelName: "test-channel",
		}
		handler := NewBudgetAlertHandler(mockSlack, cfg)

		message := PubSubMessage{
			// JSONとして不正なデータを入れる
			Data: []byte("invalid-json-data"),
		}

		// Execute
		err := handler.HandleBudgetAlert(ctx, message)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse pubsub message")
		assert.Equal(t, 0, len(mockSlack.Messages))
	})

	t.Run("Slack送信に失敗した場合、エラーが返される", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		mockSlack := slack.NewMockClient()
		mockSlack.Error = assert.AnError
		cfg := &config.Config{
			ChannelName: "test-channel",
		}
		handler := NewBudgetAlertHandler(mockSlack, cfg)

		alertData := map[string]interface{}{
			"budgetDisplayName": "テスト予算",
			"costAmount":        500.0,
			"budgetAmount":      1000.0,
			"currencyCode":      "USD",
		}
		jsonData, _ := json.Marshal(alertData)

		message := PubSubMessage{
			Data: jsonData,
		}

		// Execute
		err := handler.HandleBudgetAlert(ctx, message)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send slack notification")
	})
}

func TestNewBudgetAlertHandler(t *testing.T) {
	t.Run("ハンドラーが正しく作成される", func(t *testing.T) {
		// Setup
		mockSlack := slack.NewMockClient()
		cfg := &config.Config{
			ChannelName: "test-channel",
		}

		// Execute
		handler := NewBudgetAlertHandler(mockSlack, cfg)

		// Assert
		assert.NotNil(t, handler)
		assert.Equal(t, mockSlack, handler.slackClient)
		assert.Equal(t, cfg, handler.cfg)
	})
}

func TestProcessBudgetAlertHTTP(t *testing.T) {
	t.Run("HTTPリクエストからメッセージを解析できる", func(t *testing.T) {
		// Setup
		ctx := context.Background()

		alertData := map[string]interface{}{
			"budgetDisplayName": "テスト予算",
			"costAmount":        300.0,
			"budgetAmount":      1000.0,
			"currencyCode":      "USD",
		}
		jsonData, _ := json.Marshal(alertData)

		httpReq := HTTPRequest{
			Message: PubSubMessage{
				Data: jsonData,
			},
		}

		// ここで `requestBody` は {"message": {"data": "Base64文字列..."}} に自動でなる
		requestBody, _ := json.Marshal(httpReq)

		// 環境変数をセット（開発環境）
		t.Setenv("ENV", "dev")
		t.Setenv("SLACK_BOT_USER_OAUTH_TOKEN", "xoxb-test-token")
		t.Setenv("CHANNEL_NAME", "test-channel")

		// Execute
		err := ProcessBudgetAlertHTTP(ctx, requestBody)

		// NOTE: 実際のSlack送信は行われないが、パースは成功する
		// エラーがないか、またはSlack送信エラーであることを確認
		if err != nil {
			// Slack送信のエラーは許容（モックがないため）
			assert.Contains(t, err.Error(), "slack")
		}
	})

	t.Run("不正なJSONの場合、エラーが返される", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		invalidJSON := []byte("invalid json")

		// Execute
		err := ProcessBudgetAlertHTTP(ctx, invalidJSON)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse request")
	})
}
