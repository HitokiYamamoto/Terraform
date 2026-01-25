package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/budgetalert"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/config"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/slack"
)

// PubSubMessage はPub/Subから受け取るメッセージの構造
type PubSubMessage struct {
	// JSON Unmarshal時に自動的にBase64デコードされるため、[]byteで受ける
	Data []byte `json:"data"`
}

// BudgetAlertHandler は予算アラートを処理するハンドラー
type BudgetAlertHandler struct {
	slackClient slack.Client
	cfg         *config.Config
}

// NewBudgetAlertHandler は新しいハンドラーを作成する
func NewBudgetAlertHandler(slackClient slack.Client, cfg *config.Config) *BudgetAlertHandler {
	return &BudgetAlertHandler{
		slackClient: slackClient,
		cfg:         cfg,
	}
}

// HandleBudgetAlert は予算アラートを処理する
func (h *BudgetAlertHandler) HandleBudgetAlert(ctx context.Context, message PubSubMessage) error {
	// string変換せず、[]byteのまま渡す
	// budgetalert.ParsePubSubMessage(data []byte) に合わせる
	alert, err := budgetalert.ParsePubSubMessage(message.Data)
	if err != nil {
		return fmt.Errorf("failed to parse pubsub message: %w", err)
	}

	// Slackメッセージをフォーマット
	slackMessage := budgetalert.FormatSlackMessage(alert)

	// Slackに通知
	if err := h.slackClient.PostMessage(h.cfg.ChannelName, slackMessage); err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}

	log.Printf("Budget alert sent successfully: %s (%.2f%%)", alert.BudgetDisplayName, alert.UsagePercentage())
	return nil
}

// ProcessBudgetAlertはCloud Functions用のエントリーポイント
func ProcessBudgetAlert(ctx context.Context, m PubSubMessage) error {
	// 1. 設定を読み込む (内部で os.Getenv を実行)
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Slackクライアントを作成 (ConfigからTokenを取得)
	slackClient := slack.NewClient(cfg.SlackToken)

	// 3. ハンドラーを作成して処理
	handler := NewBudgetAlertHandler(slackClient, cfg)
	return handler.HandleBudgetAlert(ctx, m)
}

// HTTPRequest はHTTPトリガー用のリクエスト構造（テスト用）
type HTTPRequest struct {
	Message PubSubMessage `json:"message"`
}

// ProcessBudgetAlertHTTP はHTTPトリガー用のエントリーポイント（テスト/開発用）
func ProcessBudgetAlertHTTP(ctx context.Context, data []byte) error {
	var req HTTPRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("failed to parse request: %w", err)
	}

	return ProcessBudgetAlert(ctx, req.Message)
}
