module "django" {
  source                 = "./modules/repository"
  repository_name        = "Django"
  description            = "Djangoお試しリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["django", "web", "framework"]
  delete_branch_on_merge = true
  archived               = false
}
module "kubernetes" {
  source                 = "./modules/repository"
  repository_name        = "Kubernetes"
  description            = "Kubernetesお試しリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["kubernetes", "container", "orchestration"]
  delete_branch_on_merge = true
  archived               = false
}

module "github_copilot" {
  source                 = "./modules/repository"
  repository_name        = "Github-Copilot"
  description            = "GitHub Copilotお試しリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["github", "copilot", "ai"]
  delete_branch_on_merge = true
  archived               = false
}

module "golang_tutorial" {
  source                 = "./modules/repository"
  repository_name        = "Golang-Tutorial"
  description            = "Golangのチュートリアル用リポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["golang", "tutorial", "learning"]
  delete_branch_on_merge = true
  archived               = false
}

module "pyspark" {
  source                 = "./modules/repository"
  repository_name        = "PySpark"
  description            = "PySparkお試しリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["pyspark", "python", "bigdata"]
  delete_branch_on_merge = true
  archived               = false
}

module "memo_app" {
  source                 = "./modules/repository"
  repository_name        = "Memo-App"
  description            = "メモアプリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["app", "memo"]
  delete_branch_on_merge = true
  archived               = false
}

module "python_unittest" {
  source                 = "./modules/repository"
  repository_name        = "Python-UnitTest"
  description            = "Pythonユニットテストリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["python", "unittest", "testing"]
  delete_branch_on_merge = true
  archived               = false
}

module "household_account_book" {
  source                 = "./modules/repository"
  repository_name        = "Household-Account-Book"
  description            = "家計簿アプリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["app", "household", "finance"]
  delete_branch_on_merge = true
  archived               = false
}

module "python_tutorial" {
  source                 = "./modules/repository"
  repository_name        = "Python-Tutorial"
  description            = "Pythonチュートリアル"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["python", "tutorial", "learning"]
  delete_branch_on_merge = true
  archived               = false
}

module "typescript_tutorial" {
  source                 = "./modules/repository"
  repository_name        = "TypeScript-Tutorial"
  description            = "TypeScriptチュートリアル"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["typescript", "tutorial", "learning"]
  delete_branch_on_merge = true
  archived               = false
}

module "python_solid_design" {
  source                 = "./modules/repository"
  repository_name        = "Python-SOLID-Design-Principles"
  description            = "PythonのSOLID設計パターン"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["python", "solid", "design-patterns"]
  delete_branch_on_merge = true
  archived               = false
}

module "web_api" {
  source                 = "./modules/repository"
  repository_name        = "Web-API"
  description            = "Web APIリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["api", "web", "rest"]
  delete_branch_on_merge = true
  archived               = false
}

module "docker" {
  source                 = "./modules/repository"
  repository_name        = "Docker"
  description            = "Dockerお試しリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["docker", "container"]
  delete_branch_on_merge = true
  archived               = false
}

module "react_tutorial" {
  source                 = "./modules/repository"
  repository_name        = "React-Tutorial"
  description            = "Reactチュートリアル"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["react", "tutorial", "learning"]
  delete_branch_on_merge = true
  archived               = false
}

module "python_practice" {
  source                 = "./modules/repository"
  repository_name        = "Python-Practice"
  description            = "Python練習用リポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["python", "practice"]
  delete_branch_on_merge = true
  archived               = false
}

module "nodejs_exercise" {
  source                 = "./modules/repository"
  repository_name        = "Nodejs-Exercise"
  description            = "Node.js演習リポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["nodejs", "javascript", "exercise"]
  delete_branch_on_merge = true
  archived               = false
}

module "python" {
  source                 = "./modules/repository"
  repository_name        = "Python"
  description            = "Pythonリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["python"]
  delete_branch_on_merge = true
  archived               = false
}

module "bootstrap" {
  source                 = "./modules/repository"
  repository_name        = "Bootstrap"
  description            = "Bootstrapお試しリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["bootstrap", "css", "frontend"]
  delete_branch_on_merge = true
  archived               = false
}

module "gallery" {
  source                 = "./modules/repository"
  repository_name        = "Gallery"
  description            = "ギャラリーアプリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["gallery", "app"]
  delete_branch_on_merge = true
  archived               = false
}

module "sidebar" {
  source                 = "./modules/repository"
  repository_name        = "Sidebar"
  description            = "サイドバーコンポーネント"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["component", "sidebar"]
  delete_branch_on_merge = true
  archived               = false
}

module "java" {
  source                 = "./modules/repository"
  repository_name        = "Java"
  description            = "Javaリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["java"]
  delete_branch_on_merge = true
  archived               = false
}

module "ruby" {
  source                 = "./modules/repository"
  repository_name        = "Ruby"
  description            = "Rubyリポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["ruby"]
  delete_branch_on_merge = true
  archived               = false
}

module "cat_care_community" {
  source                 = "./modules/repository"
  repository_name        = "Cat-Care-Community"
  description            = "猫のヘルスケアコミュニティアプリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["cat", "community", "app"]
  delete_branch_on_merge = true
  archived               = true
}

# Organization用の特別なリポジトリを作成
module "organization_profile" {
  source                 = "./modules/repository"
  repository_name        = "HitokiYamamoto"
  description            = "Organization Profile"
  visibility             = "public"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["github", "organization"]
  delete_branch_on_merge = true
  archived               = false
}

module "cli" {
  source                 = "./modules/repository"
  repository_name        = "CLI"
  description            = "CLI作成用リポジトリ"
  visibility             = "private"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["golang", "cli", "command-line"]
  delete_branch_on_merge = true
  archived               = false
}

module "terraform" {
  source                 = "./modules/repository"
  repository_name        = "Terraform"
  description            = "インフラリソース管理用リポジトリ"
  visibility             = "public"
  has_issues             = false
  has_wiki               = false
  auto_init              = true
  topics                 = ["terraform", "infrastructure", "iac"]
  delete_branch_on_merge = true
  archived               = false
}
