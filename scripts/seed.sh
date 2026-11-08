#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

DB_DIR="${HOME}/Library/Application Support/net-cli"
DB_PATH="${DB_DIR}/net-cli.db"
SEED_FILE="${ROOT_DIR}/internal/storage/seeds/001_initial_data.sql"

if [[ ! -f "$DB_PATH" ]]; then
    echo "Database not found: $DB_PATH"
    echo "Run net-cli first to initialize the database."
    exit 1
fi

sqlite3 "$DB_PATH" < "$SEED_FILE"

echo "Database seeded successfully."