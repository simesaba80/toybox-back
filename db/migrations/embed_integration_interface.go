//go:build !integration

package migrations

import "embed"

// EmbedFiles は通常ビルド時のダミー定義です（統合テスト以外では使用しません）。
// このファイルがないとtestdb.goでエラーが発生します
var EmbedFiles embed.FS
