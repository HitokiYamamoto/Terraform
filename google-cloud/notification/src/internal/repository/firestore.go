package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

// State はFirestoreに保存する状態データ
// `firestore` タグを使ってDB側のフィールド名とマッピング
type State struct {
	LastThreshold float64   `firestore:"last_threshold"` // 前回通知したしきい値 (例: 0.5, 0.8)
	LastHeartbeat time.Time `firestore:"last_heartbeat"` // 最後に通知をした日時
	CurrentMonth  string    `firestore:"current_month"`  // 現在処理中の月 (例: "2026-02")
}

// Client はFirestore操作をラップする構造体
type Client struct {
	fs *firestore.Client
}

// NewClient はFirestoreクライアントを初期化する
func NewClient(ctx context.Context, projectID string, databaseID string) (*Client, error) {
	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		return nil, err
	}

	return &Client{
		fs: client,
	}, nil
}

// Close はクライアント接続を閉じる。
func (c *Client) Close() error {
	return c.fs.Close()
}

// SaveState は状態を保存（作成または上書き）する
func (c *Client) SaveState(ctx context.Context, docID string, state *State) error {
	// コレクション名 "billing_notifications"、ドキュメントIDごとに保存
	_, err := c.fs.Collection("billing_notifications").Doc(docID).Set(ctx, state)
	return err
}

// GetState は状態を取得する
// ドキュメントが存在しない場合、呼び出し元で codes.NotFound として判定できるエラーを返す
func (c *Client) GetState(ctx context.Context, docID string) (*State, error) {
	docSnap, err := c.fs.Collection("billing_notifications").Doc(docID).Get(ctx)
	if err != nil {
		// ドキュメントがない場合、ここでエラーが返る（status code = NotFound）
		// これをそのまま返すことで、テスト側のcodes.NotFound 判定が通る
		return nil, err
	}

	var state State
	// 取得したデータを構造体にマッピング
	if err := docSnap.DataTo(&state); err != nil {
		return nil, err
	}

	return &state, nil
}
