package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("トークンを指定した場合、Slackクライアントが正常に作成される", func(t *testing.T) {
		// Execute
		client := NewClient("test-token")

		// Assert
		assert.NotNil(t, client)
		assert.NotNil(t, client.client)
	})
}

func TestMockClient(t *testing.T) {
	t.Run("メッセージを送信した場合、モックに記録される", func(t *testing.T) {
		// Setup
		mock := NewMockClient()

		// Execute
		err := mock.PostMessage("test-channel", "テストメッセージ")

		// Assert
		assert.NoError(t, err)
		assert.Len(t, mock.Messages, 1)
		lastMsg := mock.GetLastMessage()
		assert.NotNil(t, lastMsg)
		assert.Equal(t, "test-channel", lastMsg.Channel)
		assert.Equal(t, "テストメッセージ", lastMsg.Text)
	})

	t.Run("エラーが設定されている場合、エラーが返される", func(t *testing.T) {
		// Setup
		mock := NewMockClient()
		mock.Error = assert.AnError

		// Execute
		err := mock.PostMessage("test-channel", "テストメッセージ")

		// Assert
		assert.Error(t, err)
		assert.Len(t, mock.Messages, 0)
	})

	t.Run("複数メッセージを送信した場合、全て記録される", func(t *testing.T) {
		// Setup
		mock := NewMockClient()

		// Execute
		mock.PostMessage("channel1", "メッセージ1")
		mock.PostMessage("channel2", "メッセージ2")

		// Assert
		assert.Len(t, mock.Messages, 2)
		assert.Equal(t, "メッセージ2", mock.GetLastMessage().Text)
	})

	t.Run("Resetした場合、メッセージ履歴がクリアされる", func(t *testing.T) {
		// Setup
		mock := NewMockClient()
		mock.PostMessage("test-channel", "テストメッセージ")

		// Execute
		mock.Reset()

		// Assert
		assert.Len(t, mock.Messages, 0)
		assert.Nil(t, mock.GetLastMessage())
	})
}

// NOTE: PostMessageの実際のテストはモックを使用して行っている
// 実際のSlack APIへの接続テストはインテグレーションテストで実施
