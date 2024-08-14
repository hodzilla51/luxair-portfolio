# ポートフォリオ

## 概要
このリポジトリは、私、小鍛治康人の技術ポートフォリオとして作成されたものです。このポートフォリオでは、私が過去に行ったプロジェクト、特にWeb開発に関するものを展示しております。

公開設定である本リポジトリは、セキュリティを鑑みて、プロジェクトからコードを抜粋し記載しております。ご興味を持っていただければ、コラボレーターとして各プロジェクトにご招待いたしますので、お気軽に[ご連絡](#連絡先)ください。

## プロジェクト
以下に、このポートフォリオに含まれる主要なプロジェクトをリストアップします。各プロジェクトには簡単な説明と、技術スタック、成果物へのリンクを含めています。

### 家具ブログプロジェクト（実験的）
- **説明**: このプロジェクトは、Next,jsをを用いた上で、LightHouseを中心にSEOスコアをどこまで引き上げられるかを挑戦した、実験的なものになります。「高速な動作」、「カスタム性」、「メンテナンス性」を追求し、一から開発いたしました。
- **技術スタック**: Next.js(App Router), Node.js, Vercel, MicroCMS
- **リンク**: https://kagu.monster/articles/08018
- **補足**: 全体的にまだ開発中で、トップページなどは粗末です。

### Sutututuプロジェクト
- **説明**: サーキット特化型SNSを謳ったサイトです。サーキット名とそのレイアウト名を入力すると、指定場所でのリザルト（タイムを速い順でソートしたもの）を表示いたします。二輪車、四輪車ともに対応しております。基本機能として、「記録の投稿」、「車種別での絞り込み」、「リザルト画面の表示」がございます。事前調査から、意外とまだ作られていなかったのと、個人的に欲しかったので作成いたしました。また、SNSサイトに必要な「不特定多数からのユーザー認証」、「ステートレスなAPI設計」、「サーバーの高速化・負担軽減」などの要求事項を満たすため、以下の技術スタックを利用しております。
- **技術スタック**: Next.js, Node.js, Vercel, Google Cloud Platform (何度も設計変更をしており、過去には他にもReact, Go, PostgreSQL, Apache Cassandra, Keycloak, AWS, Dockerなどを使用しておりました。）
- **リンク**: https://sutututu.vercel.app
- **広告戦略**: 過去にティッシュ配りによる広告を実施し、CVRは約3%を達成しました。
- **プロジェクトサンプル**:
サンプルとして、プロジェクトの進行のため作成した設計書と、バックエンド関連のコードをご用意させていただきました。下記はルーティングのコードです。

- [プロジェクト設計書集](./project2/docs/) - プロジェクトの全体的な設計書とドキュメント。
- [インターフェース層ルーター](./project2/samples-backend/interface/routers/) - バックエンドのルーティングロジック。

また、下記はバックエンド処理の始点と終点です。全体の流れをご覧いただくため、ご用意いたしました。お時間がありましたら、ご確認ください。

- [インフラ層ミドルウェア](./project2/samples-backend/infrastructure/middleware/) - システムのミドルウェア構成。
- [インターフェース層リポジトリ](./project2/samples-backend/interface/repositories/) - データアクセス層のリポジトリコード。
- [main.go](./project2/samples-backend/main.go) - goのエントリポイント。
- [Dockerfile](./project2/samples-backend/Dockerfile) - 簡単なものですが、Dockerfile。

- **補足**: 現在、Flutterを用いたアプリ版を開発中です。GPSを用いたラップタイムの自動測定機能を追加し、投稿画面におけるユーザーの入力項目を減らし、UXを向上させる努力を続けております。

## スキル
ここでは、私が持っている技術的なスキルセットを挙げさせていただきます。
以下の技術に知見があります：
- フロントエンド開発: HTML, CSS, JavaScript, TypeScript ,React, Next.js(App Router)
- バックエンド開発: Go, node.js, PHP
- データベース: PostgreSQL, Cassandra
- 生成AI: GPT-4, DALL-E 3
- その他: Docker, AWS, Keycloak, Git

Keycloakの設定を一から行ったことで、セキュリティに対する基礎知識を習得いたしました。また、日常的にWSL、Dockerなどの仮想Linux環境に触れておりますので、Linuxに関する基本的な操作・構造の理解をしております。
ちなみに、最もできると自認しているのはNext.js(TypeScript)です。

## 連絡先
何か気になる点がございましたら、以下の方法でお気軽にご連絡ください：
- **Email**: [lemontea.craft@gmail.com](mailto:lemontea.craft@gmail.com)
- **Tell**: 090-1813-9275

# Lisence
This project is licensed under the MIT License, see the LICENSE.txt file for details
