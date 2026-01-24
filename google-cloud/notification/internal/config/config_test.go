package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("開発環境の場合、環境変数から設定を読み込む", func(t *testing.T) {
		// Setup
		os.Setenv("ENV", "dev")
		os.Setenv("SLACK_BOT_USER_OAUTH_TOKEN", "xoxb-test-token")
		os.Setenv("CHANNEL_NAME", "test-channel")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("SLACK_BOT_USER_OAUTH_TOKEN")
			os.Unsetenv("CHANNEL_NAME")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "dev", cfg.Env)
		assert.True(t, cfg.IsDev())
		assert.Equal(t, "xoxb-test-token", cfg.SlackToken)
		assert.Equal(t, "test-channel", cfg.ChannelName)
	})

	t.Run("本番環境の場合、Google CloudプロジェクトIDを設定する", func(t *testing.T) {
		// Setup
		os.Setenv("ENV", "prod")
		os.Setenv("GOOGLE_CLOUD_PROJECT_ID", "test-project")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("GOOGLE_CLOUD_PROJECT_ID")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "prod", cfg.Env)
		assert.False(t, cfg.IsDev())
		assert.Equal(t, "test-project", cfg.ProjectID)
	})

	t.Run("ENV変数が未設定の場合、デフォルトで開発環境になる", func(t *testing.T) {
		// Setup
		os.Unsetenv("ENV")
		os.Setenv("SLACK_BOT_USER_OAUTH_TOKEN", "token")
		os.Setenv("CHANNEL_NAME", "channel")
		defer func() {
			os.Unsetenv("SLACK_BOT_USER_OAUTH_TOKEN")
			os.Unsetenv("CHANNEL_NAME")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "dev", cfg.Env)
		assert.True(t, cfg.IsDev())
	})

	t.Run("開発環境でSLACK_BOT_USER_OAUTH_TOKENが未設定の場合、エラーが返される", func(t *testing.T) {
		// Setup
		os.Setenv("ENV", "dev")
		os.Unsetenv("SLACK_BOT_USER_OAUTH_TOKEN")
		os.Setenv("CHANNEL_NAME", "test-channel")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("CHANNEL_NAME")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "SLACK_BOT_USER_OAUTH_TOKEN")
	})

	t.Run("開発環境でCHANNEL_NAMEが未設定の場合、エラーが返される", func(t *testing.T) {
		// Setup
		os.Setenv("ENV", "dev")
		os.Setenv("SLACK_BOT_USER_OAUTH_TOKEN", "xoxb-test-token")
		os.Unsetenv("CHANNEL_NAME")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("SLACK_BOT_USER_OAUTH_TOKEN")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "CHANNEL_NAME")
	})

	t.Run("本番環境でGOOGLE_CLOUD_PROJECT_IDが未設定の場合、エラーが返される", func(t *testing.T) {
		// Setup
		os.Setenv("ENV", "prod")
		os.Unsetenv("GOOGLE_CLOUD_PROJECT_ID")
		defer func() {
			os.Unsetenv("ENV")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "GOOGLE_CLOUD_PROJECT_ID")
	})
}
