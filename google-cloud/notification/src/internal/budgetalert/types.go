package budgetalert

// BudgetAlert は予算アラートの情報を保持する
type BudgetAlert struct {
	BudgetDisplayName string  // 予算名
	AlertThreshold    float64 // アラート閾値（例: 0.5 = 50%）
	CostAmount        float64 // 現在のコスト
	BudgetAmount      float64 // 予算額
	CurrencyCode      string  // 通貨コード（USD, JPYなど）
	CostIntervalStart string  // 例: "2026-01-01T00:00:00Z"
}

// UsagePercentage は予算の使用率を計算する（0-100の範囲）
func (b *BudgetAlert) UsagePercentage() float64 {
	if b.BudgetAmount == 0 {
		return 0
	}
	return (b.CostAmount / b.BudgetAmount) * 100
}

// IsOverBudget は予算を超過しているかを返す
func (b *BudgetAlert) IsOverBudget() bool {
	return b.CostAmount > b.BudgetAmount
}

// AlertLevel はアラートのレベルを返す
func (b *BudgetAlert) AlertLevel() string {
	percentage := b.UsagePercentage()

	if percentage >= 100 {
		return "🚨 危険"
	} else if percentage >= 80 {
		return "⚠️ 警告"
	} else if percentage >= 50 {
		return "ℹ️ 注意"
	}
	return "✅ 正常"
}
