package budgetalert

import (
	"encoding/base64"
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
		base64Data := base64.StdEncoding.EncodeToString(jsonData)

		// Execute
		alert, err := ParsePubSubMessage(base64Data)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, alert)
		assert.Equal(t, "テスト予算", alert.BudgetDisplayName)
		assert.Equal(t, 0.5, alert.AlertThreshold)
		assert.Equal(t, 500.0, alert.CostAmount)
		assert.Equal(t, 1000.0, alert.BudgetAmount)
		assert.Equal(t, "USD", alert.CurrencyCode)
	})

	t.Run("Base64デコードに失敗した場合、エラーが返される", func(t *testing.T) {
		// Setup: 不正なBase64文字列
		invalidBase64 := "これは不正なBase64です!!!"

		// Execute
		alert, err := ParsePubSubMessage(invalidBase64)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, alert)
		assert.Contains(t, err.Error(), "failed to decode base64")
	})

	t.Run("JSONパースに失敗した場合、エラーが返される", func(t *testing.T) {
		// Setup: 不正なJSONをBase64エンコード
		invalidJSON := "これは不正なJSONです"
		base64Data := base64.StdEncoding.EncodeToString([]byte(invalidJSON))

		// Execute
		alert, err := ParsePubSubMessage(base64Data)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, alert)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("空の文字列の場合、エラーが返される", func(t *testing.T) {
		// Execute
		alert, err := ParsePubSubMessage("")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, alert)
	})

	t.Run("必須フィールドが欠けている場合でも、パース可能な部分は取得できる", func(t *testing.T) {
		// Setup: 一部のフィールドのみのJSONデータ
		alertData := map[string]interface{}{
			"budgetDisplayName": "部分データ予算",
			"costAmount":        300.0,
		}
		jsonData, _ := json.Marshal(alertData)
		base64Data := base64.StdEncoding.EncodeToString(jsonData)

		// Execute
		alert, err := ParsePubSubMessage(base64Data)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, alert)
		assert.Equal(t, "部分データ予算", alert.BudgetDisplayName)
		assert.Equal(t, 300.0, alert.CostAmount)
		assert.Equal(t, 0.0, alert.BudgetAmount) // デフォルト値
	})
}
