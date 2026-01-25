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

// テストで使用する固定のドキュメントID
const testDocID = "test-budget-item"

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

// テスト実行ごとにDBをきれいにする
func (s *FirestoreTestSuite) clearDB() {
	// ★修正: 未定義だった docID の代わりに、定数 testDocID を使用
	_, err := s.client.fs.Collection("billing_notifications").Doc(testDocID).Delete(s.ctx)
	// ドキュメントが存在しない場合の削除エラーは無視して良いが、
	// Firestoreのエミュレータの挙動としてDeleteは成功扱いになるはず
	if err != nil && status.Code(err) != codes.NotFound {
		s.Require().NoError(err)
	}
}

func (s *FirestoreTestSuite) TestRepository_Scenarios() {
	fixedTime := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		setup     func()
		act       func() (*State, error)
		assertion func(got *State, err error)
	}{
		{
			name: "【正常系】データが何もない場合、SaveStateすると、正しく保存される",
			setup: func() {
				s.clearDB()
			},
			act: func() (*State, error) {
				newState := &State{
					LastThreshold: 0.5,
					CurrentMonth:  "2026-02",
					LastHeartbeat: fixedTime,
				}

				err := s.client.SaveState(s.ctx, testDocID, newState)
				if err != nil {
					return nil, err
				}

				return s.client.GetState(s.ctx, testDocID)
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
				oldState := &State{LastThreshold: 0.2, CurrentMonth: "2026-01"}

				s.client.SaveState(s.ctx, testDocID, oldState)
			},
			act: func() (*State, error) {
				newState := &State{LastThreshold: 0.8, CurrentMonth: "2026-02"}

				err := s.client.SaveState(s.ctx, testDocID, newState)
				if err != nil {
					return nil, err
				}

				return s.client.GetState(s.ctx, testDocID)
			},
			assertion: func(got *State, err error) {
				s.Require().NoError(err)
				s.Assert().Equal(0.8, got.LastThreshold, "値が0.2から0.8に更新されていること")
			},
		},
		{
			name: "【異常系】データが存在しない場合、GetStateすると、NotFoundエラーになる",
			setup: func() {
				s.clearDB()
			},
			act: func() (*State, error) {

				return s.client.GetState(s.ctx, testDocID)
			},
			assertion: func(got *State, err error) {
				s.Assert().Error(err)
				s.Assert().Equal(codes.NotFound, status.Code(err), "エラーコードがNotFoundであること")
				s.Assert().Nil(got)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setup != nil {
				tt.setup()
			}
			got, err := tt.act()
			if tt.assertion != nil {
				tt.assertion(got, err)
			}
		})
	}
}

func TestFirestoreTestSuite(t *testing.T) {
	suite.Run(t, new(FirestoreTestSuite))
}
