package slack

// MockClient はテスト用のモックSlackクライアント
type MockClient struct {
	Messages []Message
	Error    error
}

// Message は送信されたメッセージを記録します
type Message struct {
	Channel string
	Text    string
}

// NewMockClient は新しいモックSlackクライアントを作成します
func NewMockClient() *MockClient {
	return &MockClient{
		Messages: make([]Message, 0),
	}
}

// PostMessage はモックでメッセージを記録します
func (m *MockClient) PostMessage(channel, text string) error {
	if m.Error != nil {
		return m.Error
	}

	m.Messages = append(m.Messages, Message{
		Channel: channel,
		Text:    text,
	})
	return nil
}

// GetLastMessage は最後に送信されたメッセージを返します
func (m *MockClient) GetLastMessage() *Message {
	if len(m.Messages) == 0 {
		return nil
	}
	return &m.Messages[len(m.Messages)-1]
}

// Reset はメッセージ履歴をクリアします
func (m *MockClient) Reset() {
	m.Messages = make([]Message, 0)
	m.Error = nil
}
