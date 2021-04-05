#!/bin/bash -e

echo "HIT ENTRYPOINT"

export GF_SERVER_HTTP_PORT="${PORT:-3000}"

if [[ "$DATABASE_URL" =~ ^postgres://([^:]+):([^@]+)@([^:]+):([^/]+)/(.*)$ ]]; then
    echo "Postgres Database Url found..."

    export GF_DATABASE_TYPE="postgres"
    export GF_DATABASE_SSL_MODE="require"
    export GF_DATABASE_HOST=${BASH_REMATCH[3]}:${BASH_REMATCH[4]}
    export GF_DATABASE_NAME=${BASH_REMATCH[5]}
    export GF_DATABASE_USER=${BASH_REMATCH[1]}
    export GF_DATABASE_PASSWORD=${BASH_REMATCH[2]}
fi

exec /run.sh
