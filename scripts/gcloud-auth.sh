#!bin/bash

# GOOGLE_CLOUD_PROJECT_IDは.envファイルに設定している
check_gcloud_auth() {
    echo "Checking gcloud authentication..."
    if ! gcloud auth print-access-token &> /dev/null; then
        echo "Authentication is invalid, logging in...";
        gcloud auth application-default login;
        gcloud config set project $GOOGLE_CLOUD_PROJECT_ID;
    else
        echo "gcloud is authenticated."
    fi
}

check_gcloud_auth
