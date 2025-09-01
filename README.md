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
