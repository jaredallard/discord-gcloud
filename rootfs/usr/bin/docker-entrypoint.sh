#!/usr/bin/env bash
# vim: set ft=sh :
set -e -o pipefail

if [[ -n "$GOOGLE_SERVICE_ACCOUNT_KEY_FILE" ]]; then
  if [[ -z "$GCLOUD_PROJECT" ]]; then
    echo "[docker-entrypoint] GCLOUD_PROJECT is not set" >&2
    exit 1
  fi

  if [[ -z "$GOOGLE_SERVICE_ACCOUNT_EMAIL" ]]; then
    echo "[docker-entrypoint] GOOGLE_SERVICE_ACCOUNT_EMAIL is not set" >&2
    exit 1
  fi

  echo "[docker-entrypoint] GOOGLE_SERVICE_ACCOUNT_KEY_FILE is set, writing to gcp.json"
  base64 -d <<<"$GOOGLE_SERVICE_ACCOUNT_KEY_FILE" >gcp.json

  echo "[docker-entrypoint] Activating service account"
  gcloud auth activate-service-account "$GOOGLE_SERVICE_ACCOUNT_EMAIL" --key-file=gcp.json --project "$GCLOUD_PROJECT"
fi

echo "[docker-entrypoint] Starting dgcloud"
exec dgcloud
