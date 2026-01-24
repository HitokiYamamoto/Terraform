package main

import (
	"context"
	"log"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/config"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/secret"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/slack"
)

func main() {
	ctx := context.Background()

	// 設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var token, channelName string

	if cfg.IsDev() {
		// 開発環境: 設定から直接取得
		token = cfg.SlackToken
		channelName = cfg.ChannelName
	} else {
		// 本番環境: Secret Managerから取得
		secretMgr, err := secret.NewManager(ctx, cfg.ProjectID)
		if err != nil {
			log.Fatalf("Failed to create secret manager: %v", err)
		}
		defer secretMgr.Close()

		token, err = secretMgr.GetSecret(ctx, "SLACK_BOT_USER_OAUTH_TOKEN")
		if err != nil {
			log.Fatalf("Failed to get slack token: %v", err)
		}

		channelName, err = secretMgr.GetSecret(ctx, "CHANNEL_NAME")
		if err != nil {
			log.Fatalf("Failed to get channel name: %v", err)
		}
	}

	// Slackクライアントを作成してメッセージを送信
	slackClient := slack.NewClient(token)
	if err := slackClient.PostMessage(channelName, "テスト通知"); err != nil {
		log.Fatalf("Failed to post message: %v", err)
	}

	log.Println("Message sent successfully!")
}
