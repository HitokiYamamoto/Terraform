package budgetalert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSlackMessage(t *testing.T) {
	t.Run("予算アラートの場合、フォーマットされたメッセージが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			BudgetDisplayName: "本番環境予算",
			AlertThreshold:    0.8,
			CostAmount:        850.50,
			BudgetAmount:      1000.00,
			CurrencyCode:      "USD",
		}

		message := FormatSlackMessage(alert)

		assert.Contains(t, message, "本番環境予算")
		assert.Contains(t, message, "850.50")
		assert.Contains(t, message, "1000.00")
		assert.Contains(t, message, "USD")
		assert.Contains(t, message, "85.05%")
		assert.Contains(t, message, "⚠️ 警告")
	})

	t.Run("予算超過の場合、危険レベルのメッセージが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			BudgetDisplayName: "開発環境予算",
			AlertThreshold:    0.9,
			CostAmount:        1200.00,
			BudgetAmount:      1000.00,
			CurrencyCode:      "JPY",
		}

		message := FormatSlackMessage(alert)

		assert.Contains(t, message, "開発環境予算")
		assert.Contains(t, message, "🚨 危険")
		assert.Contains(t, message, "120.00%")
	})

	t.Run("通貨コードがJPYの場合、円記号が含まれる", func(t *testing.T) {
		alert := &BudgetAlert{
			BudgetDisplayName: "テスト予算",
			CostAmount:        50000,
			BudgetAmount:      100000,
			CurrencyCode:      "JPY",
		}

		message := FormatSlackMessage(alert)

		assert.Contains(t, message, "JPY")
		assert.Contains(t, message, "50000")
	})

	t.Run("使用率が低い場合、正常レベルのメッセージが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			BudgetDisplayName: "少額予算",
			CostAmount:        30,
			BudgetAmount:      100,
			CurrencyCode:      "USD",
		}

		message := FormatSlackMessage(alert)

		assert.Contains(t, message, "✅ 正常")
		assert.Contains(t, message, "30.00%")
	})
}
