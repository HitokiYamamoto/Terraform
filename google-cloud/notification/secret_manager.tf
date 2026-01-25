module "slack_oauth_token" {
  source    = "../modules/secret_manager"
  secret_id = "slack-bot-user-oauth-token"
}

module "slack_channel_name" {
  source    = "../modules/secret_manager"
  secret_id = "slack-channel-name"
}
