#!/bin/bash

set -e

app="$1"
org="$2"

if [ -z "$app" ]; then
    echo "Usage: $0 <heroku_app_name>"
    exit 1
fi

apps=$(heroku apps --org "$org")

if ! echo "${apps}" | grep "^$app$" >/dev/null; then
    echo "$app not found..."
    heroku create "$app" --org "$org"
fi

app_info=$(heroku apps:info --app "$app")

# Setting stack to container
if ! echo "${app_info}" | grep "^Stack:.*container$" >/dev/null; then
    heroku stack:set container --app "$app"
fi

app_addons=$(heroku addons --app "$app")

# Creating heroku-postgresql:hobby-dev
if ! echo "${app_addons}" | grep "^heroku-postgresql" >/dev/null; then
    heroku addons:create heroku-postgresql:hobby-dev --app "$app"
fi
