# help
[default]
help:
    @just --list

# Auth Login Google Cloud
gcloud-login:
    docker compose run --rm \
        terraform \
        bash ./scripts/gcloud-auth.sh

# terraform commands
terraform +args:
    docker compose run --rm terraform \
        terraform -chdir=bootstrap {{ args }}

# terraform commands with service account impersonation
terraform-impersonate service_account +args:
    docker compose run --rm \
        --env GOOGLE_IMPERSONATE_SERVICE_ACCOUNT={{ service_account }} \
        terraform \
        terraform -chdir=google-cloud {{ args }}

# terraform fmt
format:
    docker compose run --rm terraform terraform fmt -recursive

# terraform validate
lint:
    docker compose run --rm terraform tflint --init
    docker compose run --rm terraform tflint --recursive

[private]
TRIVY_VERSION := "0.69.0" # renovate: datasource=docker depName=aquasec/trivy
# security scan with trivy
trivy:
    docker run --rm -it \
        --mount type=bind,source=$(pwd),target=/app \
        --workdir /app \
        aquasec/trivy:{{ TRIVY_VERSION  }} \
        filesystem \
        --ignorefile .trivyignore \
        --skip-files "bootstrap/main.tf" \
        --scanners misconfig,vuln .

[private]
PRETTIER_VERSION := "3.8.1" # renovate: datasource=npm depName=prettier
# Prettier for YAML and JSON5 files
prettier:
    docker compose run --rm node \
        npx prettier@{{ PRETTIER_VERSION }} --write "**/*.yaml" "**/*.json5"

[private]
YAMLLINT_VERSION := "1.35.1" # renovate: datasource=github-releases depName=adrienverge/yamllint
# Lint check for YAML files
yamllint:
    docker compose run --rm python \
        uv tool run yamllint@{{ YAMLLINT_VERSION  }} .

# Golangのタスク実行(default: task list)
task +args="":
    docker compose run --rm golang \
        task {{ args }}

# Golangのユニットテスト実行
test:
    docker compose up firestore --detach
    docker compose run --rm golang task test
    docker compose down

# 予算アラートのテストメッセージをPub/Subに発行
budget-publish-test:
    docker compose run --rm terraform \
        gcloud pubsub topics publish budget-alert-topic \
        --message='{ \
            "budgetDisplayName": "テスト予算アラート", \
            "alertThresholdExceeded": 0.0, \
            "costAmount": 0, \
            "budgetAmount": 10000, \
            "currencyCode": "JPY", \
        }'

[private]
RENOVATE_VERSION := "42.94.6"  # renovate: datasource=docker depName=renovate/renovate
# Renovateの設定ファイル検証
renovate-validate:
    @docker run --rm -it \
        --mount type=bind,source="$(pwd)/.github/renovate.json5",target=/target/renovate.json5 \
        --workdir /target \
        renovate/renovate:{{ RENOVATE_VERSION  }} \
        renovate-config-validator

# Renovateのデバッグ実行
renovate-debug:
    @docker run --rm -it \
        --env RENOVATE_TOKEN=${GITHUB_TOKEN} \
        --env LOG_LEVEL=debug \
        --env RENOVATE_CONFIG_FILE="/target/.github/renovate.json5" \
        --env RENOVATE_BASE_BRANCHES=`git branch --show-current` \
        --env LOG_FORMAT=json \
        --mount type=bind,source="$(pwd)/.github/renovate.json5",target=/target/.github/renovate.json5 \
        renovate/renovate:{{ RENOVATE_VERSION  }} \
        --require-config=ignored \
        --dry-run=full \
        "HitokiYamamoto/Terraform" > .vscode/renovate.log 2>&1

    @cat .vscode/renovate.log | jq -r ' \
    select(.branchesInformation) \
    | .branchesInformation[] \
    | .branchName as $branch \
    | .result as $result \
    | .upgrades[] \
    | [$branch, $result, .depName, .packageFile] \
    | @tsv' | column -t -s $'\t'
