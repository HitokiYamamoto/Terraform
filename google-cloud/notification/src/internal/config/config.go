package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーションの設定を保持する
type Config struct {
	SlackToken  string
	ChannelName string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	_ = godotenv.Load(".env") // エラーは無視（環境変数から読み込む）

	cfg := &Config{
		// 修正点: 「=」を「:」に変更し、行末に「,」を追加
		SlackToken:  os.Getenv("SLACK_BOT_USER_OAUTH_TOKEN"),
		ChannelName: os.Getenv("CHANNEL_NAME"),
	}

	if cfg.SlackToken == "" {
		return nil, fmt.Errorf("環境変数 SLACK_BOT_USER_OAUTH_TOKEN が設定されていません")
	}
	if cfg.ChannelName == "" {
		return nil, fmt.Errorf("環境変数 CHANNEL_NAME が設定されていません")
	}

	return cfg, nil
}
