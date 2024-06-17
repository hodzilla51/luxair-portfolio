# Cassandraテーブル定義設計書

## テーブル名: `posts`

### 概要
user_id、およびbucketをパーティションキーとした投稿用テーブル。
特定ユーザーの投稿に特化。
データの特性から、多様なクエリが考えられるためPotgresに移行。

### Primary Key
- Partition Key: user_id, bucket
- Clustering Columns: post_id(DESC)

### Columns
| Column Name | Data Type | Description |
|-------------|-----------|-------------|
| user_id   | text    | UUID        　　|
| bucket    | int     | パーティション用 |
| post_id   | bigint  | Snowflake使用   |
| post_type | text    | fuel, eventなど |
| vehicle_id| text    | UUID            |
| reply_to  | bigint  | post_id         |
| repost_to | bigint  | post_id         |
| text      | text    | 本文            |
| image_folder_url | text | AWSフォルダ|

### Indexes
- `index_name` on `column_name`: Indexの説明。

### 制約
- ALLOW FILTERING禁止！！
- パーティション一つあたり10MBあたりを超える場合、bucketを追加

### クエリパターン
特定のユーザーの最新の投稿を取得するクエリ：

```cql
SELECT * FROM posts WHERE user_id = '特定のuser_id' AND bucket = 特定のbucket番号 ORDER BY post_id DESC LIMIT 10;
```

前回取得時、以降に作成された新しい投稿に限定して取得するクエリ：
※DESCのため古い順の取得は不可

```cql
SELECT * FROM posts
WHERE user_id = ? AND bucket = ? AND post_id > ?
ORDER BY post_id DESC
LIMIT 10;
```

### 使用中のテーブル定義
```cql
CREATE TABLE posts (
  user_id text,
  bucket int,
  post_id bigint,
  post_type text,
  vehicle_id text,
  reply_to bigint,
  repost_to bigint,
  text text,
  image_folder_url text,
  PRIMARY KEY ((user_id, bucket), post_id)
) WITH CLUSTERING ORDER BY (post_id DESC);
```

## テーブル名: `posts_month`

### 概要
posted_week、およびbucketをパーティションキーとした投稿用テーブル。
投稿データ全体の取得に特化。週ごとに別パーティションとし、サイズが大きくなりすぎた場合bucketで区切る。
容量、一貫性の確保の難度から将来的に削除を予定。

### Primary Key
- Partition Key: posted_month, bucket
- Clustering Columns: post_id(DESC)

### Columns
| Column Name | Data Type | Description |
|-------------|-----------|-------------|
| posted_month | text    | RFC3339を月区切り|
| bucket       | int     | パーティション用 |
| post_id      | bigint  | Snowflake使用   |
| post_type    | text    | fuel, eventなど |
| user_id      | text    | UUID        　　|
| vehicle_id   | text    | UUID            |
| reply_to     | bigint  | post_id         |
| repost_to    | bigint  | post_id         |
| text         | text    | 本文            |
| image_folder_url | text | AWSフォルダ|

### Indexes
- `index_name` on `column_name`: Indexの説明。

### 制約
- ALLOW FILTERING禁止！！
- 実質投稿データが二倍となるため将来的に削除を検討。

### クエリパターン
特定のユーザーの最新の投稿を取得するクエリ：

```cql
SELECT * FROM posts WHERE user_id = '特定のuser_id' AND bucket = 特定のbucket番号 ORDER BY post_id DESC LIMIT 10;
```

前回取得時、以降に作成された新しい投稿に限定して取得するクエリ：
※DESCのため古い順の取得は不可

```cql
SELECT * FROM posts_month
WHERE posted_month = '2023-03' AND bucket = 1
ORDER BY post_id DESC
LIMIT 10;
```

### 使用中のテーブル定義
```cql
CREATE TABLE posts_month (
  posted_month text,
  bucket int,
  post_id bigint,
  post_type text,
  user_id text,
  vehicle_id text,
  reply_to bigint,
  repost_to bigint,
  text text,
  image_folder_url text,
  PRIMARY KEY ((posted_month, bucket), post_id)
) WITH CLUSTERING ORDER BY (post_id DESC);
```

## テーブル名: `fuels`

### 概要
vehicle_id, bucketをパーティションキーとしたガソリン記録用テーブル。
特定乗り物に対してのリクエストに特化。

### Primary Key
- Partition Key: vehicle_id, bucket
- Clustering Columns: fuel_id(DESC)

### Columns
| Column Name | Data Type | Description |
|-------------|-----------|-------------|
| vehicle_id      | text    | UUID         |
| bucket          | int     | パーティション用 |
| fuel_id         | bigint  | Snowflake使用   |
| post_id         | bigint  | Snowflake使用   |
| refual_at       |timestamp| RFC3339       |
| amount_ml       | int     | 給油量(単位ml) |
| fuel_type       | text    | 油種。dieselとか |
| liter_fee       | int     | リッター当たりの値段 |
| total_fee       | int     | 今回給油にかかった金額 |
| odometer_kilo   | int     | 給油時点の走行距離 |
| tripmeter_kilo  | int     | 前回の給油からの距離 |
| location        | text    | 場所 |
| receipt_img_url | text    | レシート画像(AWS) |


### Indexes
- `index_name` on `column_name`: Indexの説明。

### 制約
- ALLOW FILTERING禁止！！

### クエリパターン
特定車両の最新の給油記録を取得するクエリ：

```cql
SELECT * FROM fuels
WHERE vehicle_id = ? AND bucket = ?
ORDER BY fuel_id DESC
LIMIT 10;
```

特定の期間の給油情報を取得するクエリ：

```cql
SELECT * FROM fuels
WHERE vehicle_id = ? AND bucket = ? AND fuel_id > ? AND fuel_id < ?
ORDER BY fuel_id DESC;
```

### 使用中のテーブル定義
```cql
CREATE TABLE fuels (
  vehicle_id text,
  bucket int,
  fuel_id bigint,
  post_id bigint,
  refuel_at timestamp,
  amount_ml int,
  fuel_type text,
  liter_fee int,
  total_fee int,
  odometer_kilo int,
  tripmeter_kilo int,
  location text,
  receipt_img_url text,
  PRIMARY KEY ((vehicle_id, bucket), refuel_at)
) WITH CLUSTERING ORDER BY (refuel_at DESC);
```

## テーブル名: `events`

### 概要
user_id、およびbucketをパーティションキーとしたイベント管理用テーブル。
これは一覧として取得するからCassandra向きやで～

### Primary Key
- Partition Key: user_id, bucket
- Clustering Columns: event_id(DESC)

### Columns
| Column Name | Data Type | Description |
|-------------|-----------|-------------|
| user_id     | text      | UUID                 |
| bucket      | int       | パーティション用      |
| event_id    | bigint    | Snowflake           |
| post_id     | text      | 実質詳細              |
| start_at    | timestamp | RFC3339             |
| end_at      | timestamp | RFC3339             |
| event_title | text      | 短く　　             |
| is_allday   | boolean   | 終日かどうか          |
| location    | text      | 場所情報。変更するかも |

### Indexes
- `index_name` on `column_name`: Indexの説明。

### 制約
- ALLOW FILTERING禁止！！
- パーティション一つあたり10MBあたりを超える場合、bucketを追加

### クエリパターン
特定のユーザー(user_id)に対するイベントの一覧を取得：

```cql
SELECT event_id, start_at, end_at, event_title, is_allday, location, post_id
FROM events
WHERE user_id = '特定のuser_id' AND bucket = 特定のバケット番号
ORDER BY event_id DESC;
```

### 使用中のテーブル定義
```cql
CREATE TABLE events (
  user_id text,
  bucket int,
  event_id bigint,
  post_id bigint,
  start_at timestamp,
  end_at timestamp,
  event_title text,
  is_allday boolean,
  location text,
  PRIMARY KEY ((user_id, bucket), event_id)
) WITH CLUSTERING ORDER BY (event_id DESC);
```

## テーブル名: `likes`

### 概要
partition_seed、およびbucketをパーティションキーとしたいいね用テーブル。
いいね機能のデータ特性から、データサイズが安定しないのがデメリット。
あと、パーティションが大量に増えるからPostgreSQLとどちらにするか検討中。いや、多分Postgresが最適解。

### Primary Key
- Partition Key: post_id, bucket
- Clustering Columns: liked_at(DESC)

### Columns
| Column Name | Data Type | Description |
|-------------|-----------|-------------|
| post_id        | bigint    | Snowflake使用   |
| bucket         | int       | パーティション用 |
| user_id        | text      | UUID           |
| liked_at       | timestamp | RFC3339        |

### Indexes
- `index_name` on `column_name`: Indexの説明。

### 制約
- ALLOW FILTERING禁止！！
- パーティション一つあたり10MBあたりを超える場合、bucketを追加

### クエリパターン
特定の投稿(post_id)に対するいいねを取得：

```cql
SELECT * FROM likes WHERE user_id = ? AND bucket = ? AND post_id = ?;
```

### 使用中のテーブル定義
```cql
CREATE TABLE likes (
    post_id bigint,
    bucket int,
    user_id text,
    liked_at timestamp,
    PRIMARY KEY ((user_id, bucket), post_id)
) WITH CLUSTERING ORDER BY (post_id DESC);
```


## テーブル名: `followings`

### 概要
user_id、およびbucketをパーティションキーとしたフォロー管理用テーブル。
likesと同様、パーティションが無数に生成される欠陥がある。ただ、こっちは投稿ほど数は増えない。Postgresへの移行を検討中。

### Primary Key
- Partition Key: user_id, bucket
- Clustering Columns: follow_id(DESC)

### Columns
| Column Name | Data Type | Description |
|-------------|-----------|-------------|
| user_id        | text      | UUID           |
| bucket         | int       | パーティション用 |
| following_id   | text      | Snowflake使用   |
| following_at   | timestamp | RFC3339        |

### Indexes
- `index_name` on `column_name`: Indexの説明。

### 制約
- ALLOW FILTERING禁止！！
- パーティション一つあたり10MBあたりを超える場合、bucketを追加

### クエリパターン
特定のユーザー(user_id)に対するフォローの一覧を取得：

```cql
SELECT following_id
FROM followings
WHERE user_id = '特定のuser_id' AND bucket = 特定のバケット番号
ORDER BY following_id DESC;
```

### 使用中のテーブル定義
```cql
CREATE TABLE following (
    user_id text,
    bucket int,
    following_id text,
    following_at
    PRIMARY KEY ((user_id, bucket), following_at)
) WITH CLUSTERING ORDER BY (following_at DESC);
```

## テーブル名: `followers`

### 概要
user_id、およびbucketをパーティションキーとしたフォロー管理用テーブル。
likesと同様、パーティションが無数に生成される欠陥がある。ただ、こっちは投稿ほど数は増えない。Postgresへの移行を検討中。

### Primary Key
- Partition Key: user_id, bucket
- Clustering Columns: followed_id(DESC)

### Columns
| Column Name | Data Type | Description |
|-------------|-----------|-------------|
| user_id        | text      | UUID           |
| bucket         | int       | パーティション用 |
| follower_id   | text      | Snowflake使用   |
| follower_at   | timestamp | RFC3339        |

### Indexes
- `index_name` on `column_name`: Indexの説明。

### 制約
- ALLOW FILTERING禁止！！
- パーティション一つあたり10MBあたりを超える場合、bucketを追加

### クエリパターン
特定のユーザー(user_id)に対するフォロワーの一覧を取得：

```cql
SELECT follower_id, follower_at
FROM followers
WHERE user_id = '特定のuser_id' AND bucket = 特定のバケット番号
ORDER BY follower_id DESC;
```

### 使用中のテーブル定義
```cql
CREATE TABLE following (
    user_id text,
    bucket int,
    follower_id text,
    follower_at
    PRIMARY KEY ((user_id, bucket), follower_at)
) WITH CLUSTERING ORDER BY (follower_at DESC);
```

