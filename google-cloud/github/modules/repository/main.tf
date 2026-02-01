# リポジトリ設定
resource "github_repository" "main" {
  name                   = var.repository_name
  description            = var.description
  visibility             = var.visibility
  has_issues             = var.has_issues
  has_wiki               = var.has_wiki
  auto_init              = var.auto_init
  topics                 = var.topics
  delete_branch_on_merge = var.delete_branch_on_merge
  vulnerability_alerts   = true
  archived               = var.archived

  /* NOTE: デフォルトブランチをmainに設定
  auto_init = true の場合、自動的にmainブランチが作成される
  */
}

/* NOTE: ブランチ保護ルールの設定

ブランチ保護ルールの設定を追加したいが、Privateリポジトリはできない
回避するには、Publicリポジトリにするか、有料プランにするしかない
*/
