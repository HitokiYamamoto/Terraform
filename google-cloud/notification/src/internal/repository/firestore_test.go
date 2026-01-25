package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreTestSuite struct {
	suite.Suite
	client *Client
	ctx    context.Context
}

func (s *FirestoreTestSuite) SetupSuite() {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		s.T().Skip("Skipping integration test: FIRESTORE_EMULATOR_HOST not set")
	}
	s.ctx = context.Background()
	var err error
	s.client, err = NewClient(s.ctx, "test-project-id", "(default)")
	s.Require().NoError(err)
}

func (s *FirestoreTestSuite) TearDownSuite() {
	if s.client != nil {
		s.client.fs.Close()
	}
}

// テスト実行ごとにDBをきれいにする（テーブル駆動ループ内の干渉を防ぐため）
func (s *FirestoreTestSuite) clearDB() {
	_, err := s.client.fs.Collection("billing_notifications").Doc("state").Delete(s.ctx)
	s.Require().NoError(err)
}

func (s *FirestoreTestSuite) TestRepository_Scenarios() {
	// 固定の日時（テスト用）
	fixedTime := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name string // テストケース名

		// Given: テスト実行前のDBの状態を作る関数
		setup func()

		// When: 検証したい操作を実行する関数 (戻り値やエラーを返す)
		act func() (*State, error)

		// Then: アサーション (結果の検証)
		assertion func(got *State, err error)
	}{
		{
			name: "【正常系】データが何もない場合、SaveStateすると、正しく保存される",
			setup: func() {
				s.clearDB() // 空にする
			},
			act: func() (*State, error) {
				// 保存したいデータ
				newState := &State{
					LastThreshold: 0.5,
					CurrentMonth:  "2026-02",
					LastHeartbeat: fixedTime,
				}
				err := s.client.SaveState(s.ctx, newState)
				if err != nil {
					return nil, err
				}
				// 検証のためにGetして返す
				return s.client.GetState(s.ctx)
			},
			assertion: func(got *State, err error) {
				s.Require().NoError(err)
				s.Assert().Equal(0.5, got.LastThreshold)
				s.Assert().Equal("2026-02", got.CurrentMonth)
				s.Assert().WithinDuration(fixedTime, got.LastHeartbeat, time.Second)
			},
		},
		{
			name: "【正常系】既にデータがある場合、SaveStateで上書きされる",
			setup: func() {
				s.clearDB()
				// 事前に古いデータを入れておく
				oldState := &State{LastThreshold: 0.2, CurrentMonth: "2026-01"}
				s.client.SaveState(s.ctx, oldState)
			},
			act: func() (*State, error) {
				// 新しいデータで上書き
				newState := &State{LastThreshold: 0.8, CurrentMonth: "2026-02"}
				err := s.client.SaveState(s.ctx, newState)
				if err != nil {
					return nil, err
				}
				return s.client.GetState(s.ctx)
			},
			assertion: func(got *State, err error) {
				s.Require().NoError(err)
				s.Assert().Equal(0.8, got.LastThreshold, "値が0.2から0.8に更新されていること")
			},
		},
		{
			name: "【異常系】データが存在しない場合、GetStateすると、NotFoundエラーになる",
			setup: func() {
				s.clearDB() // 完全に空にする
			},
			act: func() (*State, error) {
				return s.client.GetState(s.ctx)
			},
			assertion: func(got *State, err error) {
				s.Assert().Error(err)
				s.Assert().Equal(codes.NotFound, status.Code(err), "エラーコードがNotFoundであること")
				s.Assert().Nil(got)
			},
		},
	}

	// テーブル駆動実行ループ
	for _, tt := range tests {
		s.Run(tt.name, func() {
			// 1. Given
			if tt.setup != nil {
				tt.setup()
			}

			// 2. When
			got, err := tt.act()

			// 3. Then
			if tt.assertion != nil {
				tt.assertion(got, err)
			}
		})
	}
}

func TestFirestoreTestSuite(t *testing.T) {
	suite.Run(t, new(FirestoreTestSuite))
}
