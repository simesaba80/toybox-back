# チャレキャラBackend  

## プロジェクト概要
部員の作成した作品やブログの投稿や閲覧ができるWebアプリ、Toyboxのバックエンドです。  
既に動いているWebアプリのリプレイスを目指すものになります。(https://github.com/Kyutech-C3/toybox-server/tree/develop)

## 開発環境の準備
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

# データの挿入(初回のみ)
$  docker compose exec -T db psql -U postgres -d toybox < toybox_backup.sql

# 以下のURLを開き、起動しているか確認
http://localhost:8080/

#サーバーの停止
$ docker compose down

```

※依存先が増えた場合は適切に/internal/di/wire.goに記述し、`wire ./internal/di`を実行すること

### internal内の各ディレクトリの役割
このプロジェクトではクリーンアーキテクチャを採用しています。各レイヤーの役割は以下のようになっています。

`/domain`
他の層に依存しない中心、Entityではこのアプリが取り扱う領域や概念の定義、repositoryではデータアクセスの抽象化を行っている。EntityのオブジェクトはDBのテーブル定義とは異なるので注意

`/usecase`
domainに依存しつつ、アプリケーション固有のビジネスロジックを記述する。Echoなどのフレームワークには依存しない

`/interface`
infrastructure層とusecase層の橋渡しを行う。  
controllerではフレームワークで受け取ったHTTPリクエストの処理、schemaではリクエストおよびレスポンスの形式を変換している

`infrastructure`
DBやフロントエンドといった外部のシステムと実際に通信する部分になる。  
repositoryでは/domain/repositoryで定義したデータアクセスの具体的な実装、routerではEchoを使ったルーティングを実際に行っている。
