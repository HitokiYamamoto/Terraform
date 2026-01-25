package budgetalert

import (
	"fmt"
)

// FormatSlackMessage は予算アラートをSlack用のメッセージにフォーマットする
func FormatSlackMessage(alert *BudgetAlert) string {
	return fmt.Sprintf(
		`📊 *予算アラート通知*

*予算名:* %s
*アラートレベル:* %s

*使用状況:*
• 現在のコスト: %.2f %s
• 予算額: %.2f %s
• 使用率: %.2f%%

%s`,
		alert.BudgetDisplayName,
		alert.AlertLevel(),
		alert.CostAmount,
		alert.CurrencyCode,
		alert.BudgetAmount,
		alert.CurrencyCode,
		alert.UsagePercentage(),
		getActionMessage(alert),
	)
}

// getActionMessage はアラートレベルに応じたアクションメッセージを返す
func getActionMessage(alert *BudgetAlert) string {
	percentage := alert.UsagePercentage()

	if percentage >= 100 {
		return "⚠️ *予算を超過しています！至急確認してください。*"
	} else if percentage >= 80 {
		return "⚠️ 予算の80%を超えています。使用状況を確認してください。"
	} else if percentage >= 50 {
		return "ℹ️ 予算の50%を超えています。今後の使用状況にご注意ください。"
	}
	return "✅ 現在の使用状況は問題ありません。"
}
