package secret

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockManager(t *testing.T) {
	t.Run("シークレットが登録されている場合、正しい値が返される", func(t *testing.T) {
		// Setup
		mock := NewMockManager()
		mock.Secrets["test-secret"] = "test-value"
		ctx := context.Background()

		// Execute
		result, err := mock.GetSecret(ctx, "test-secret")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "test-value", result)
	})

	t.Run("シークレットが登録されていない場合、空文字列が返される", func(t *testing.T) {
		// Setup
		mock := NewMockManager()
		ctx := context.Background()

		// Execute
		result, err := mock.GetSecret(ctx, "non-existent")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("エラーが設定されている場合、エラーが返される", func(t *testing.T) {
		// Setup
		mock := NewMockManager()
		mock.Error = errors.New("test error")
		ctx := context.Background()

		// Execute
		result, err := mock.GetSecret(ctx, "test-secret")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, "", result)
	})

	t.Run("複数のシークレットを登録した場合、それぞれ正しく取得できる", func(t *testing.T) {
		// Setup
		mock := NewMockManager()
		mock.Secrets["secret1"] = "value1"
		mock.Secrets["secret2"] = "value2"
		mock.Secrets["secret3"] = "value3"
		ctx := context.Background()

		// Execute & Assert
		val1, err1 := mock.GetSecret(ctx, "secret1")
		assert.NoError(t, err1)
		assert.Equal(t, "value1", val1)

		val2, err2 := mock.GetSecret(ctx, "secret2")
		assert.NoError(t, err2)
		assert.Equal(t, "value2", val2)

		val3, err3 := mock.GetSecret(ctx, "secret3")
		assert.NoError(t, err3)
		assert.Equal(t, "value3", val3)
	})

	t.Run("NewMockManagerで作成した場合、Secretsマップが初期化される", func(t *testing.T) {
		// Execute
		mock := NewMockManager()

		// Assert
		assert.NotNil(t, mock)
		assert.NotNil(t, mock.Secrets)
		assert.Equal(t, 0, len(mock.Secrets))
	})
}
