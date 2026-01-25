// Package notification はGoogle Cloud Functions用のエントリーポイント
package notification

import (
	"context"

	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/function"
)

// ProcessBudgetAlert はCloud Functions用のエントリーポイント
func ProcessBudgetAlert(ctx context.Context, m function.PubSubMessage) error {
	return function.ProcessBudgetAlert(ctx, m)
}
