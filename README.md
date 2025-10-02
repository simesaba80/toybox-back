# チャレキャラ Backend

## プロジェクト概要

部員の作成した作品やブログの投稿や閲覧ができる Web アプリ、Toybox のバックエンドです。  
既に動いている Web アプリのリプレイスを目指すものになります。(https://github.com/Kyutech-C3/toybox-server)
/docs/design-doc.md に[設計ドキュメント](docs/design-doc.md)があります。

## 開発環境の準備

### 開発に必要な CLI ツールのインストール

以下のコマンド実行後、wire, migrate, swag コマンドがそれぞれ実行できるようになることを確認してください。実行できない場合 `go env GOPATH`で表示されるディレクトリのパスが通ってない可能性が高いです。

```
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@lates
$ go install github.com/google/wire/cmd/wire@latest
$ go install github.com/swaggo/swag/cmd/swag@latest
```

### プロジェクトのセットアップ

```
# レポジトリをクローンする
$ git clone git@github.com:simesaba80/toybox-back.git

# 作業ディレクトリに移動する
$ cd toybox-back

# 依存モジュールをインストールする
$ go mod download

# .envの作成と書き込み
$ cp .env.example .env

# サーバーを起動
$ docker compose up -d

# マイグレーションの実行
$ migrate --path db/migrations --database '<DBの接続文字列>' up

# データの挿入(初回のみ)
# バックアップのSQLをプロジェクト直下に配置後
$ ./scripts/movedata.sh

# 以下のURLを開き、起動しているか確認
http://localhost:8080/

#サーバーの停止
$ docker compose down

```

※依存先が増えた場合は適切に/internal/di/wire.go に記述し、`wire ./internal/di`を実行すること

### データベースのセットアップ

データベースのマイグレーションには golang-migrate を使っています。

基本的な使い方は以下の記事を参考にしてください。
https://zenn.dev/farstep/books/f74e6b76ea7456/viewer/4cd440
マイグレーション時 DSN の指定が必要になります。コンテナのネットワーク外から実行するので.env で記入した DSN と HOST の値が異なることに注意してください

### API ドキュメントの更新

http://localhost:8080/swagger/index.html にアクセスすることで API ドキュメントを確認できる。

API ドキュメントの更新

```
$ swag init -g cmd/main.go --parseDependency
```
