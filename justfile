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

# security scan with trivy
trivy:
    docker run --rm -it \
        --mount type=bind,source=$(pwd),target=/app \
        --workdir /app \
        aquasec/trivy:0.68.2 fs --scanners config,vuln .

# Prettier for YAML and JSON5 files
prettier:
    docker compose run --rm prettier --write "**/*.yaml" "**/*.json5"

# Lint check for YAML files
yamllint:
    docker compose run --rm yamllint .

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
