# GitHubリポジトリ管理

## 使用方法

### リポジトリ作成

main.tfに下記を追加

```terraform
module "golang_tutorial" {
  source                 = "./modules/github"                    # moduleの場所
  repository_name        = "golang-tutorial"                     # リポジトリ名
  description            = "Golangのチュートリアル用リポジトリ"       # リポジトリの説明
  has_issues             = true                                  # issuesを有効にするかどうか
  has_wiki               = false                                 # wikiを有効にするかどうか
  auto_init              = true                                  # READMEを自動作成するかどうか
  topics                 = ["golang", "tutorial", "learning"]    # リポジトリのタグ
  delete_branch_on_merge = true                                  # マージ後にブランチを自動削除するかどうか
}
```

リポジトリがすでにある場合は、直下にimport.tfを作成して、以下を記載

```terraform
import {
  to = module.github.repository.golang_tutorial.github_repository.main  # importするmoduleの名前
  id = "golang-tutorial"                              # GitHubリポジトリの名前（backend.tfにorgを記載してるのでここではorgを省略可）
}
```

## トークンについて

1. [トークン一覧](https://github.com/settings/tokens)にアクセス
2. 今のトークンを削除（または新しく作成）
3. "Generate new token (classic)"をクリック
4. 必要なスコープをチェック
   - ✅repo (全部)
   - ✅ delete_repo(destroyする時に多分いる)
5. トークンを生成してコピー
6. 環境変数を更新: `export GITHUB_TOKEN=ghp_コピーしたトークン`

⚠️ 永続的にしたかったら以下のようにできる

ただ、トークンにはexpireを設定するはずなのでここはよしなにどうぞって感じ

```shell
echo 'export GITHUB_TOKEN=ghp_xxxxxxxxxxxxx' >> ~/.zshrc
source ~/.zshrc
```
