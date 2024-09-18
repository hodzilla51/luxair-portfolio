# 概要
このリポジトリは、私、小鍛治康人の技術ポートフォリオとして作成されたものです。このポートフォリオでは、私が過去に行ったプロジェクト、特にWeb開発に関するものを展示しております。

公開設定である本リポジトリは、セキュリティを鑑みて、プロジェクトからコードを抜粋し記載しております。ご興味を持っていただければ、コラボレーターとして各プロジェクトにご招待いたしますので、お気軽に[ご連絡](#連絡先)ください。

# プロジェクト
以下に、このポートフォリオに含まれる主要なプロジェクトをリストアップします。各プロジェクトには簡単な説明と、技術スタック、成果物へのリンクを含めています。

## Sutututuプロジェクト
[Sutututuを見てみる](https://sutututu.net) <br/>
サーキット特化型SNSを謳ったサイトです。サーキット名とそのレイアウト名を入力すると、指定場所でのリザルト（タイムを速い順でソートしたもの）を表示いたします。二輪車、四輪車ともに対応しております。
### **説明**
基本機能として、「記録の投稿」、「車種別での絞り込み」、「リザルト画面の表示」がございます。事前調査から、意外とまだ作られていなかったのと、個人的に欲しかったので作成いたしました。また、SNSサイトに必要な「不特定多数からのユーザー認証」、「ステートレスなAPI設計」、「サーバーの高速化・負担軽減」などの要求事項を満たすために、以下に記載の技術を使用しております。
### **技術スタック**
Next.js, Node.js, Vercel, Google Cloud Platform (何度も設計変更をしており、過去には他にもReact(生), Go, PostgreSQL, Apache Cassandra, Keycloak, AWS, Dockerなどを使用しておりました。）
### 技術的なこだわり
 - GSAPライブラリを用いたアニメーション付きUI←視覚的に楽しいサイトを心掛け、軽量かつ高機能なGSAPライブラリを至る所で使用しております。
 - 用途毎のデータベースの変更によるFirebaseへの負荷軽減←例えば、検索補助機能はpublic内にJSONを配置し、データベースとしています。数KB程度なのでパフォーマンスに問題はありません。
 - ちょっと古いブラウザへの互換性←全ての機能がiPhoneXでも問題なく動作するように設計しております。例えば、Next.jsのuseSearchParams関数はiPhoneXでは動作しませんので、回避しています。
### **リンク**
https://sutututu.net
### **広告戦略**
過去にティッシュ配りによる広告活動を実施し、CVR（ユーザー登録 per ティッシュ）は約3%を記録しました。
### **プロジェクトサンプル**
現在稼働中のSutututuに関しては、セキュリティの懸念があるのでコードは控えさせていただきます。
代わりに、旧Sutututuのソースコードの一部を公開いたします。同プロジェクトは、ReactやGoなどを用いており、Kubernetes上で動作させることを前提に開発しました。
ほぼ完成しておりましたが、運営費の見積もりが甘く資金上の理由から頓挫した背景がございます。
- [プロジェクト設計書集](./project2/docs/) - プロジェクトの全体的な設計書とドキュメント。
- [インターフェース層ルーター](./project2/samples-backend/interface/routers/) - バックエンドのルーティングロジック。

また、下記はバックエンド処理の始点と終点です。全体の流れをご覧いただくため、ご用意いたしました。お時間がありましたら、ご確認ください。

- [インフラ層ミドルウェア](./project2/samples-backend/infrastructure/middleware/) - システムのミドルウェア構成。
- [インターフェース層リポジトリ](./project2/samples-backend/interface/repositories/) - データアクセス層のリポジトリコード。
- [main.go](./project2/samples-backend/main.go) - goのエントリポイント。
- [Dockerfile](./project2/samples-backend/Dockerfile) - 簡単なものですが、Dockerfile。


## 家具ブログプロジェクト（実験的）
[カグモンを見てみる](https://kagu.monster/articles/08018) <br/>
家具についての特化ブログサイトです。このプロジェクトでは、Next.jsをを用いた上でLightHouseを中心にSEOスコアをどこまで引き上げられるかの挑戦をしております。現在はまともに定期更新をしておらず、力試し的な意味合いの強いサイトです。
### **説明**
「高速な動作」、「カスタム性」、「メンテナンス性」を追求し、一から開発いたしました。SSG, ISRを含むサーバーサイドレンダリングの多用や、軽量ライブラリ・自作の軽量コンポーネントの使用により、LightHouseスコアを超高得点（基本的にすべて100点）で安定させております。また、適切なメタタグやサイトマップの設定、URL構造により、権威性を除くサイト設計上の基礎的なSEOに関しても最適化された状態にあります。
### **技術スタック**
Next.js(App Router), Node.js, Vercel, MicroCMS
### **リンク**
https://kagu.monster/articles/08018
### **補足**
全体的にまだ開発中で、トップページなどは未実装です。


# スキル
ここでは、私が持っている技術的なスキルセットを挙げさせていただきます。
以下の技術に知見があります。
- フロントエンド開発: HTML, CSS, JavaScript, TypeScript ,React, Next.js(App Router)
- バックエンド開発: Go, node.js, PHP
- データベース: PostgreSQL, Cassandra
- クラウド: AWS, GCP(Firebase)
- 生成AI: GPT-4, o1preview, DALL-E 3
- その他: Docker, Keycloak, Git, Vercel

また、過去にWordPressサイトを企画から運営まで一人で行い、月間30,000pvまで伸ばした経験がございます。
この際、SEO関連の知識も培っております。先述のSutututuに適用し、現在の主なサイト訪問者は検索流入からによるものとなっています。

Keycloakの設定を一から行ったことで、主に認証関連のセキュリティに対する基礎知識を習得しております。また、日常的にWSL、Dockerなどの仮想Linux環境に触れておりますので、Linuxに関する基本的な操作・構造の理解をしております。
ちなみに、最もできると自認しているのはNext.js(TypeScript)です。

# 連絡先
何か気になる点がございましたら、以下の方法でお気軽にご連絡ください：
- **Email**: [lemontea.craft@gmail.com](mailto:lemontea.craft@gmail.com)
- **Tell**: 090-1813-9275

# License
This project is licensed under the MIT License, see the LICENSE.txt file for details
