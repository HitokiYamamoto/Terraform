package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/budgetalert"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/config"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/secret"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/slack"
)

// PubSubMessage はPub/Subから受け取るメッセージの構造
type PubSubMessage struct {
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
	// Pub/Subメッセージをパース
	alert, err := budgetalert.ParsePubSubMessage(string(message.Data))
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

// ProcessBudgetAlert はCloud Functions用のエントリーポイント
func ProcessBudgetAlert(ctx context.Context, m PubSubMessage) error {
	// 設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var token string

	if cfg.IsDev() {
		// 開発環境: 設定から直接取得
		token = cfg.SlackToken
	} else {
		// 本番環境: Secret Managerから取得
		secretMgr, err := secret.NewManager(ctx, cfg.ProjectID)
		if err != nil {
			return fmt.Errorf("failed to create secret manager: %w", err)
		}
		defer secretMgr.Close()

		token, err = secretMgr.GetSecret(ctx, "SLACK_BOT_USER_OAUTH_TOKEN")
		if err != nil {
			return fmt.Errorf("failed to get slack token: %w", err)
		}

		channelName, err := secretMgr.GetSecret(ctx, "CHANNEL_NAME")
		if err != nil {
			return fmt.Errorf("failed to get channel name: %w", err)
		}
		cfg.ChannelName = channelName
	}

	// Slackクライアントを作成
	slackClient := slack.NewClient(token)

	// ハンドラーを作成して処理
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
