package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーションの設定を保持する
type Config struct {
	Env         string
	ProjectID   string
	SlackToken  string
	ChannelName string
}

// Load は環境に応じて設定を読み込む
func Load() (*Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	cfg := &Config{
		Env: env,
	}

	if env == "dev" {
		// 開発環境: .envファイルから読み込み（ファイルがなければ環境変数を使用）
		_ = godotenv.Load(".env") // エラーは無視（環境変数から読み込む）
		cfg.SlackToken = os.Getenv("SLACK_BOT_USER_OAUTH_TOKEN")
		cfg.ChannelName = os.Getenv("CHANNEL_NAME")

		// 開発環境の必須環境変数をチェック
		if cfg.SlackToken == "" {
			return nil, fmt.Errorf("環境変数 SLACK_BOT_USER_OAUTH_TOKEN が設定されていません")
		}
		if cfg.ChannelName == "" {
			return nil, fmt.Errorf("環境変数 CHANNEL_NAME が設定されていません")
		}
	} else {
		// 本番環境: Secret Managerから取得するため、ProjectIDのみ設定
		cfg.ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT_ID")

		// 本番環境の必須環境変数をチェック
		if cfg.ProjectID == "" {
			return nil, fmt.Errorf("環境変数 GOOGLE_CLOUD_PROJECT_ID が設定されていません（本番環境では必須です）")
		}
	}

	return cfg, nil
}

// IsDev は開発環境かどうかを返します
func (c *Config) IsDev() bool {
	return c.Env == "dev"
}
