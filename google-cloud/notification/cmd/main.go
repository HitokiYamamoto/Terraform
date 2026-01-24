package main

import (
	"log"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/config"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/slack"
)

func main() {

	// 設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var token, channelName string
	token = cfg.SlackToken
	channelName = cfg.ChannelName

	// Slackクライアントを作成してメッセージを送信
	slackClient := slack.NewClient(token)
	if err := slackClient.PostMessage(channelName, "テスト通知"); err != nil {
		log.Fatalf("Failed to post message: %v", err)
	}

	log.Println("Message sent successfully!")
}
