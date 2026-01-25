package budgetalert

import (
	"encoding/json"
	"fmt"
)

// pubSubMessageData はPub/Subメッセージのデータ構造
type pubSubMessageData struct {
	BudgetDisplayName      string  `json:"budgetDisplayName"`
	AlertThresholdExceeded float64 `json:"alertThresholdExceeded"`
	CostAmount             float64 `json:"costAmount"`
	BudgetAmount           float64 `json:"budgetAmount"`
	CurrencyCode           string  `json:"currencyCode"`
}

// ParsePubSubMessage はPub/Subメッセージ(JSONバイト列)をパースする
func ParsePubSubMessage(data []byte) (*BudgetAlert, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty message data")
	}

	// JSONパース
	var msgData pubSubMessageData
	if err := json.Unmarshal(data, &msgData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// BudgetAlertに変換
	alert := &BudgetAlert{
		BudgetDisplayName: msgData.BudgetDisplayName,
		AlertThreshold:    msgData.AlertThresholdExceeded,
		CostAmount:        msgData.CostAmount,
		BudgetAmount:      msgData.BudgetAmount,
		CurrencyCode:      msgData.CurrencyCode,
	}

	return alert, nil
}
