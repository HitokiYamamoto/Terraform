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

module "terraform" {
  source                 = "./modules/repository"
  repository_name        = "Terraform"
  description            = "インフラリソース管理用リポジトリ"
  visibility             = "public"
  has_issues             = true
  has_wiki               = false
  auto_init              = true
  topics                 = ["terraform", "infrastructure", "iac"]
  delete_branch_on_merge = true
  archived               = false
}

module "zenn" {
  source                 = "./modules/repository"
  repository_name        = "Zenn"
  description            = "Zenn記事管理用リポジトリ"
  visibility             = "private"
  has_issues             = true
  has_wiki               = false
  auto_init              = true
  topics                 = ["zenn", "articles", "blog"]
  delete_branch_on_merge = true
  archived               = false
}
