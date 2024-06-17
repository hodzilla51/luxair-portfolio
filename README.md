# Luxiar様向け ポートフォリオ

## 概要
このリポジトリは、私、小鍛治康人の技術ポートフォリオとしてLuxiar様向けに特別に作成されたものです。このポートフォリオでは、私が過去に行ったプロジェクト、特にWeb開発に関するものを展示しております。

公開設定である本リポジトリは、セキュリティを鑑みて、プロジェクトからコードを抜粋し記載しております。ご興味を持っていただければ、コラボレーターとしてプロジェクトにご招待いたしますので、お気軽に[ご連絡](#連絡先)ください。

## プロジェクト
以下に、このポートフォリオに含まれる主要なプロジェクトをリストアップします。各プロジェクトには簡単な説明と、技術スタック、成果物へのリンクを含めています。

### プロジェクト1
- **説明**: このプロジェクトは、新しいブログをはじめるにあたり、「高速な動作」、「カスタム性」、「メンテナンス性」が要求されたため、WordPressではなく一から開発されました。特に、Next.jsを利用して、SEO、ユーザー体験、を向上させます。
- **技術スタック**: Next.js(App Router), Node.js, Vercel, MicroCMS
- **リンク**: https://kagu.monster/articles/08018
- **補足**: 全体的にまだ開発中で、トップページなどは粗末です。

### プロジェクト2
- **説明**: 車愛好家のためのSNSサイトです。SNSサイトに必要な「不特定多数からのユーザー認証」、「RestAPIによるステートレス化」、「サーバーの高速化・負担軽減」などの事項を満たすため、以下の技術スタックを利用しております。
- **技術スタック**: React, Go, PostgreSQL, Apache Cassandra, Keycloak, AWS, Docker
- **プロジェクトサンプル**:
サンプルとして、プロジェクトの進行のため作成した設計書と、バックエンド関連のコードをご用意させていただきました。下記はルーティングのコードです。

- [プロジェクト設計書集](./project2/docs/) - プロジェクトの全体的な設計書とドキュメント。
- [インターフェース層ルーター](./project2/samples-backend/interface/routers/) - バックエンドのルーティングロジック。

また、下記はバックエンド処理の始点と終点です。全体の流れをご覧いただくため、ご用意いたしました。お時間がありましたら、ご確認ください。

- [インフラ層ミドルウェア](./project2/samples-backend/infrastructure/middleware/) - システムのミドルウェア構成。
- [インターフェース層リポジトリ](./project2/samples-backend/interface/repositories/) - データアクセス層のリポジトリコード。
- [main.go](./project2/samples-backend/main.go) - goのエントリポイント。
- [Dockerfile](./project2/samples-backend/Dockerfile) - 簡単なものですが、Dockerfile。

- **補足**: 開発中です。

## スキル
ここでは、私が持っている技術的なスキルセットを挙げさせていただきます。
以下の技術に知見があります：
- フロントエンド開発: HTML, CSS, JavaScript, TypeScript ,React, Next.js(App Router)
- バックエンド開発: Go, (一部、node.js, PHP)
- データベース: PostgreSQL, Cassandra
- その他: Docker, AWS, Keycloak

ちなみに、最もできると自認しているのはGoです。興味もバックエンド開発にあります。

## 連絡先
何か気になる点がございましたら、以下の方法でお気軽にご連絡ください：
- **Email**: [lemontea.craft@gmail.com](mailto:lemontea.craft@gmail.com)
- **Tell**: 090-1813-9275

# Lisence
This project is licensed under the MIT License, see the LICENSE.txt file for details