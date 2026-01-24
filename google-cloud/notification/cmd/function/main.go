// Package main はGoogle Cloud Functions用のエントリーポイント
package main

import (
	"context"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/HitokiYamamoto/Terraform/google-cloud/notification/internal/function"
)

func main() {
	// Cloud Functions用に登録
	funcframework.RegisterEventFunction("/", processBudgetAlert)

	// ローカル開発用のHTTPサーバーを起動
	port := "8080"
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("Failed to start functions framework: %v", err)
	}
}

// init はCloud Functionsにデプロイする際の初期化関数
func init() {
	// Cloud Functions Gen 1用に登録
	funcframework.RegisterEventFunction("/", processBudgetAlert)
}

// processBudgetAlert はCloud Functions用のエントリーポイント
func processBudgetAlert(ctx context.Context, m function.PubSubMessage) error {
	return function.ProcessBudgetAlert(ctx, m)
}
