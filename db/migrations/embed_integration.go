//go:build integration

package migrations

import "embed"

// Files は統合テスト用にマイグレーションSQLを埋め込んだ仮想ファイルシステムです。
//
//go:embed *.up.sql
var EmbedFiles embed.FS
