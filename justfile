# help
[default]
help:
    @just --list --unsorted

# justfileのフォーマット
fmt:
    @just --fmt --unstable

# Auth Login Google Cloud
gcloud-login:
    docker compose run --rm \
        terraform \
        bash ./scripts/gcloud-auth.sh

# terraform commands
terraform-bootstrap +args:
    docker compose run --rm terraform \
        terraform -chdir=google-cloud/bootstrap {{ args }}

# terraform commands with service account impersonation
terraform-impersonate service_account +args:
    docker compose run --rm \
        --env GOOGLE_IMPERSONATE_SERVICE_ACCOUNT={{ service_account }} \
        terraform \
        terraform -chdir=google-cloud {{ args }}

# terraform fmt
terraform-format:
    docker compose run --rm terraform terraform fmt -recursive

# terraform validate && terraform tflint
terraform-lint:
    docker compose run --rm terraform tflint --init && tflint --recursive
    docker compose run --rm terraform tflint --recursive
    docker compose run --rm terraform terraform -chdir=google-cloud validate

[private]
TRIVY_VERSION := "0.69.0@sha256:33f816d414b9d582d25bb737ffa4a632ae34e222f7ec1b50252cef0ce2266006"

# security scan with trivy
trivy:
    docker run --rm -it \
        --mount type=bind,source=$(pwd),target=/app \
        --workdir /app \
        aquasec/trivy:{{ TRIVY_VERSION }} \
        filesystem \
        --ignorefile .trivyignore \
        --skip-files "google-cloud/bootstrap/main.tf" \
        --scanners misconfig,vuln .

[private]
PRETTIER_VERSION := "3.8.1"

# Prettier for YAML and JSON5 files
prettier:
    docker compose run --rm node \
        npx prettier@{{ PRETTIER_VERSION }} --write "**/*.yaml" "**/*.json5"

[private]
YAMLLINT_VERSION := "v1.38.0"

# Lint check for YAML files
yamllint:
    docker compose run --rm python \
        uv tool run yamllint@{{ YAMLLINT_VERSION }} .

# Golangのタスク実行(default: task list)
task +args="":
    #!/usr/bin/env zsh
    # testはfirestoreエミュレータ起動後に実行するため除外
    if [ "{{ args }}" = "test" ]; then
        just task-test
    else
        docker compose run --rm golang task {{ args }}
    fi

# Golangのユニットテスト実行
task-test:
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
RENOVATE_VERSION := "43.0.7@sha256:bc63974fbbf5505779cc937b3d7f8858c6bad9815cc04de17b18b37966e51c6a"

# Renovateの設定ファイル検証
renovate-validate:
    @docker run --rm -it \
        --mount type=bind,source="$(pwd)/.github/renovate.json5",target=/target/renovate.json5 \
        --workdir /target \
        renovate/renovate:{{ RENOVATE_VERSION }} \
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
        renovate/renovate:{{ RENOVATE_VERSION }} \
        --require-config=ignored \
        --dry-run=full \
        "HitokiYamamoto/Terraform" > .vscode/renovate.log 2>&1

    @cat .vscode/renovate.log | jq -r -R ' \
    fromjson? \
    | select(.branchesInformation) \
    | .branchesInformation[] \
    | .branchName as $branch \
    | .result as $result \
    | .upgrades[] \
    | [$branch, $result, .depName, .packageFile] \
    | @tsv' | column -t -s $'\t'
