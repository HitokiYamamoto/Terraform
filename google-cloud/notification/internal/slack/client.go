package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// Client はSlackへメッセージを送信するインターフェース
type Client interface {
	PostMessage(channel, text string) error
}

// SlackClient はSlackのクライアント
type SlackClient struct {
	client *slack.Client
}

// NewClient は新しいSlackクライアントを作成する
func NewClient(token string) *SlackClient {
	return &SlackClient{
		client: slack.New(token),
	}
}

// PostMessage はSlackにメッセージを投稿する
func (s *SlackClient) PostMessage(channel, text string) error {
	/* NOTE:
	MsgOptionText()の第二引数は、メッセージ内でMarkdownを使用するかどうかを指定する。
	falseに設定すると、プレーンテキストとして扱われる。
	*/
	_, _, err := s.client.PostMessage(
		channel,
		slack.MsgOptionText(text, false),
	)
	if err != nil {
		return fmt.Errorf("failed to post message to slack: %w", err)
	}
	return nil
}
