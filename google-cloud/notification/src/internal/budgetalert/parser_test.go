package budgetalert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePubSubMessage(t *testing.T) {
	t.Run("正しいPub/Subメッセージの場合、BudgetAlertが返される", func(t *testing.T) {
		// Setup: 予算アラートのJSONデータを作成
		alertData := map[string]interface{}{
			"budgetDisplayName":      "テスト予算",
			"alertThresholdExceeded": 0.5,
			"costAmount":             500.0,
			"budgetAmount":           1000.0,
			"currencyCode":           "USD",
		}
		jsonData, _ := json.Marshal(alertData)

		// Execute
		alert, err := ParsePubSubMessage(jsonData)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, alert)
		assert.Equal(t, "テスト予算", alert.BudgetDisplayName)
		assert.Equal(t, 0.5, alert.AlertThreshold)
		assert.Equal(t, 500.0, alert.CostAmount)
		assert.Equal(t, 1000.0, alert.BudgetAmount)
		assert.Equal(t, "USD", alert.CurrencyCode)
	})

	t.Run("JSONパースに失敗した場合、エラーが返される", func(t *testing.T) {
		// Setup: 不正なJSONデータ（バイト列）
		invalidJSON := []byte("これは不正なJSONです")

		// Execute
		alert, err := ParsePubSubMessage(invalidJSON)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, alert)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("空のデータの場合、エラーが返される", func(t *testing.T) {
		// Execute
		// nilまたは空のバイトスライスを渡す
		alert, err := ParsePubSubMessage(nil)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, alert)
		assert.Contains(t, err.Error(), "empty message data")
	})

	t.Run("必須フィールドが欠けている場合でも、パース可能な部分は取得できる", func(t *testing.T) {
		// Setup: 一部のフィールドのみのJSONデータ
		alertData := map[string]interface{}{
			"budgetDisplayName": "部分データ予算",
			"costAmount":        300.0,
		}
		jsonData, _ := json.Marshal(alertData)

		// Execute
		alert, err := ParsePubSubMessage(jsonData)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, alert)
		assert.Equal(t, "部分データ予算", alert.BudgetDisplayName)
		assert.Equal(t, 300.0, alert.CostAmount)
		assert.Equal(t, 0.0, alert.BudgetAmount) // デフォルト値
	})
}
