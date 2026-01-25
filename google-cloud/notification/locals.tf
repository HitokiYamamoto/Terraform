locals {
  build_roles = [
    "roles/logging.logWriter",
    "roles/artifactregistry.writer",
    "roles/storage.objectViewer",
    "roles/cloudbuild.builds.builder",
  ]

  secrets = ({
    SLACK_BOT_USER_OAUTH_TOKEN = module.slack_oauth_token.secret_id
    CHANNEL_NAME               = module.slack_channel_name.secret_id,
  })
}
