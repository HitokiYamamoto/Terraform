package budgetalert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBudgetAlert_UsagePercentage(t *testing.T) {
	t.Run("予算の50%を使用している場合、50.0が返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   500,
			BudgetAmount: 1000,
		}

		result := alert.UsagePercentage()

		assert.Equal(t, 50.0, result)
	})

	t.Run("予算の100%を使用している場合、100.0が返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   1000,
			BudgetAmount: 1000,
		}

		result := alert.UsagePercentage()

		assert.Equal(t, 100.0, result)
	})

	t.Run("予算を超過している場合、100を超える値が返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   1500,
			BudgetAmount: 1000,
		}

		result := alert.UsagePercentage()

		assert.Equal(t, 150.0, result)
	})

	t.Run("予算が0の場合、0が返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   100,
			BudgetAmount: 0,
		}

		result := alert.UsagePercentage()

		assert.Equal(t, 0.0, result)
	})
}

func TestBudgetAlert_IsOverBudget(t *testing.T) {
	t.Run("予算を超過している場合、trueが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   1500,
			BudgetAmount: 1000,
		}

		assert.True(t, alert.IsOverBudget())
	})

	t.Run("予算内の場合、falseが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   500,
			BudgetAmount: 1000,
		}

		assert.False(t, alert.IsOverBudget())
	})

	t.Run("予算と同額の場合、falseが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   1000,
			BudgetAmount: 1000,
		}

		assert.False(t, alert.IsOverBudget())
	})
}

func TestBudgetAlert_AlertLevel(t *testing.T) {
	t.Run("使用率が50%未満の場合、正常レベルが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   400,
			BudgetAmount: 1000,
		}

		assert.Equal(t, "✅ 正常", alert.AlertLevel())
	})

	t.Run("使用率が50-80%の場合、注意レベルが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   600,
			BudgetAmount: 1000,
		}

		assert.Equal(t, "ℹ️ 注意", alert.AlertLevel())
	})

	t.Run("使用率が80-100%の場合、警告レベルが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   900,
			BudgetAmount: 1000,
		}

		assert.Equal(t, "⚠️ 警告", alert.AlertLevel())
	})

	t.Run("使用率が100%以上の場合、危険レベルが返される", func(t *testing.T) {
		alert := &BudgetAlert{
			CostAmount:   1500,
			BudgetAmount: 1000,
		}

		assert.Equal(t, "🚨 危険", alert.AlertLevel())
	})
}
