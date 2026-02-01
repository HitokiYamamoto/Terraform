/* NOTE: ブランチ保護ルールの設定

ブランチ保護ルールの設定を追加は、無料プランのPrivateリポジトリはサポートされていない。
回避するには、Publicリポジトリにするか、有料プランにするしかない
*/

resource "github_branch_protection" "main" {
  repository_id = var.repository_id

  pattern = "main"

  enforce_admins = true // 管理者にもブランチ保護を適用

  allows_force_pushes = false // mainブランチでの強制プッシュを禁止

  require_signed_commits = true // 署名されたコミットのみを許可（GitHub Appは自動的に署名される）

  required_pull_request_reviews {
    required_approving_review_count = 0    // 「レビューは不要だが、PRという形式は必須」
    dismiss_stale_reviews           = true // 新しいコミットがプッシュされた場合、既存のレビューを取り消す
  }

  // GitHub Actionsのジョブが成功していることを必須にする
  required_status_checks {
    strict   = true
    contexts = var.job_list # Actionsのジョブ名
  }
}
