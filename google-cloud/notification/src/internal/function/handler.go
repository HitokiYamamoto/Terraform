package function

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/budgetalert"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/config"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/repository"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/src/internal/slack"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PubSubMessage はPub/Subから受け取るメッセージの構造
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type StateRepository interface {
	GetState(ctx context.Context) (*repository.State, error)
	SaveState(ctx context.Context, state *repository.State) error
}

// Handler構造体の修正
type BudgetAlertHandler struct {
	slackClient slack.Client
	repo        StateRepository
	cfg         *config.Config
}

// NewBudgetAlertHandler は新しいハンドラーを作成する
func NewBudgetAlertHandler(slackClient slack.Client, repo StateRepository, cfg *config.Config) *BudgetAlertHandler {
	return &BudgetAlertHandler{
		slackClient: slackClient,
		repo:        repo,
		cfg:         cfg,
	}
}

// HandleBudgetAlert は予算アラートを処理する
func (h *BudgetAlertHandler) HandleBudgetAlert(ctx context.Context, message PubSubMessage) error {
	// Pub/Subメッセージのパース
	alert, err := budgetalert.ParsePubSubMessage(message.Data)
	if err != nil {
		return fmt.Errorf("failed to parse pubsub message: %w", err)
	}

	// 受信内容
	log.Printf("📩 予算アラートを受信しました: %+v", alert)

	// 前回の状態をFirestoreから取得
	state, err := h.repo.GetState(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// ★日本語ログ2: 初回起動時
			log.Println("前回の状態が見つかりません。新規に状態を初期化します。")
			state = &repository.State{}
		} else {
			return fmt.Errorf("failed to get state from firestore: %w", err)
		}
	}

	shouldNotify := false
	notificationNote := ""

	// --- ロジックA: 月替わりのリセット判定 ---
	if state.CurrentMonth != alert.CostIntervalStart {
		// ★日本語ログ3: 月替わり
		log.Printf("📅 月が替わりました (%s -> %s)。しきい値をリセットします。", state.CurrentMonth, alert.CostIntervalStart)
		state.LastThreshold = 0.0
		state.CurrentMonth = alert.CostIntervalStart
	}

	// --- ロジックB: しきい値上昇の判定 ---
	if alert.AlertThreshold > state.LastThreshold {
		shouldNotify = true
		state.LastThreshold = alert.AlertThreshold
	}

	// --- ロジックC: 週次生存確認 (Heartbeat) ---
	// 「現在時刻 - 前回の確認時刻」が「7日 × 24時間」を超えているか判定
	now := time.Now()
	if now.Sub(state.LastHeartbeat) > 7*24*time.Hour {
		shouldNotify = true
		notificationNote = "\n(System Heartbeat: 正常稼働中)"
		state.LastHeartbeat = now
	}

	// 通知不要ならここで終了
	if !shouldNotify {
		// スキップ理由
		log.Printf(
			"🔕 通知スキップ: 今回のしきい値(%.2f)は前回(%.2f)以下であり、月(%s)も変わっていないため。",
			alert.AlertThreshold,
			state.LastThreshold,
			state.CurrentMonth,
		)
		return nil
	}

	// Slackメッセージをフォーマット & 送信
	slackMessage := budgetalert.FormatSlackMessage(alert)

	if notificationNote != "" {
		log.Println("💓 生存確認(Heartbeat)として通知を送信します。")
	}

	if err := h.slackClient.PostMessage(h.cfg.ChannelName, slackMessage); err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}

	// 新しい状態をFirestoreに保存
	if err := h.repo.SaveState(ctx, state); err != nil {
		return fmt.Errorf("failed to save state to firestore: %w", err)
	}

	// 完了
	log.Printf("✅ 完了: Slack通知を送信し、状態を更新しました。現在の消化率: %.2f%%", alert.UsagePercentage())
	return nil
}

// ProcessBudgetAlertはCloud Functions用のエントリーポイント
func ProcessBudgetAlert(ctx context.Context, m PubSubMessage) error {
	// 設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Slackクライアントを作成
	slackClient := slack.NewClient(cfg.SlackToken)

	// Firestoreクライアントを作成
	// GOOGLE_CLOUD_PROJECT はCloud Functionsで自動的に設定される環境変数
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// DB名: (default) または 名前付きDB名
	repo, err := repository.NewClient(ctx, projectID, "(default)")
	if err != nil {
		return fmt.Errorf("failed to create firestore client: %w", err)
	}
	// 関数の終了時に接続を閉じる
	defer repo.Close()

	// ハンドラーを作成して処理 (repoを渡す)
	handler := NewBudgetAlertHandler(slackClient, repo, cfg)
	return handler.HandleBudgetAlert(ctx, m)
}
