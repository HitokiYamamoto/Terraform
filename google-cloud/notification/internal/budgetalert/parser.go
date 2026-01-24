package budgetalert

import (
	"encoding/base64"
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

// ParsePubSubMessage はBase64エンコードされたPub/Subメッセージをパースする
func ParsePubSubMessage(base64Data string) (*BudgetAlert, error) {
	if base64Data == "" {
		return nil, fmt.Errorf("empty message data")
	}

	// Base64デコード
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// JSONパース
	var data pubSubMessageData
	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// BudgetAlertに変換
	alert := &BudgetAlert{
		BudgetDisplayName: data.BudgetDisplayName,
		AlertThreshold:    data.AlertThresholdExceeded,
		CostAmount:        data.CostAmount,
		BudgetAmount:      data.BudgetAmount,
		CurrencyCode:      data.CurrencyCode,
	}

	return alert, nil
}
