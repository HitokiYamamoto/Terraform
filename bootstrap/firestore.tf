/* NOTE:
    Pub/Subへの通知は「予算のしきい値を超えた時だけ」ではなく、「1日に数回（実際には約10分〜30分おき）、現在の状況を定期的に」送信される。

    Slackへの通知が多すぎて困ったので、「Pub/Sub → Cloud Functions → Slack」の構成で、
    通知頻度を「しきい値を超えた時だけ」などに減らして、Cloud Functions内のコードでフィルタリング（条件分岐）を行う。

    Pub/Subから受け取るJSONデータ内に`alertThresholdExceeded`というフィールドがあり、
    これが直近のしきい値を超えているかを示す値（1.0, 0.9など）を持っているが、これだけだと「超えている間はずっと通知」されてしまう。

    前回通知した時のしきい値をFirestoreなどに保存しておき、今回受け取ったしきい値と比較し、
    「前回よりも高いしきい値を超えた場合のみSlackに通知する」というロジックを組むことで解決しようという戦法。

    Cloud Functions自体は1日に数十回起動し続けるが、無料枠（月間200万回呼び出し等）の中に収まる見込み。

    データの例:
    {
        "last_threshold": 0.5,           // 前回通知したしきい値 (50%)
        "last_heartbeat": "2026-01-17T...", // 最後に通知を送った日時（週に1度はパイプラインが生きていることを確認するために用意する）
        "current_month": "2026-01"       // 現在処理中の請求月
    }
*/

resource "google_firestore_database" "database" {
  project = var.project_id
  /* NOTE:
  # 名前付きデータベースは無料割り当て対象外なので、デフォルトデータベースを使用する
  # 他用途で利用したい場合は、コレクション名などで区別し、dataとして使う
  */
  name        = "(default)"
  location_id = "us-central1"
  type        = "FIRESTORE_NATIVE"

  delete_protection_state = "DELETE_PROTECTION_ENABLED"
}
