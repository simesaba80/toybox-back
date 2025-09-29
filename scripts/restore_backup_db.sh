#!/bin/bash
set -e

echo "Starting database restore..."

# データベースが既に存在するかチェック
if psql -U postgres -lqt | cut -d \| -f 1 | grep -qw toybox; then
    echo "Database toybox already exists, dropping it..."
    dropdb -U postgres toybox
fi

# データベースを作成
echo "Creating database toybox..."
createdb -U postgres toybox

# SQLファイルを実行
echo "Restoring from backup..."
psql -U postgres -d toybox -f /tmp/toybox_backup.sql

echo "Database restore completed successfully!"