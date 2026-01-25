package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("環境変数から設定を読み込むことができる", func(t *testing.T) {
		// Setup
		os.Setenv("SLACK_BOT_USER_OAUTH_TOKEN", "xoxb-test-token")
		os.Setenv("CHANNEL_NAME", "test-channel")
		defer func() {
			os.Unsetenv("SLACK_BOT_USER_OAUTH_TOKEN")
			os.Unsetenv("CHANNEL_NAME")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "xoxb-test-token", cfg.SlackToken)
		assert.Equal(t, "test-channel", cfg.ChannelName)
	})

	t.Run("SLACK_BOT_USER_OAUTH_TOKENが未設定の場合、エラーが返される", func(t *testing.T) {
		// Setup
		os.Unsetenv("SLACK_BOT_USER_OAUTH_TOKEN")
		os.Setenv("CHANNEL_NAME", "test-channel")
		defer func() {
			os.Unsetenv("CHANNEL_NAME")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "SLACK_BOT_USER_OAUTH_TOKEN")
	})

	t.Run("CHANNEL_NAMEが未設定の場合、エラーが返される", func(t *testing.T) {
		// Setup
		os.Setenv("SLACK_BOT_USER_OAUTH_TOKEN", "xoxb-test-token")
		os.Unsetenv("CHANNEL_NAME")
		defer func() {
			os.Unsetenv("SLACK_BOT_USER_OAUTH_TOKEN")
		}()

		// Execute
		cfg, err := Load()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "CHANNEL_NAME")
	})
}
