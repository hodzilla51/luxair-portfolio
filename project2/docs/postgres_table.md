# Postgresテーブル定義設計書

## テーブル名: `posts`

### 概要
`user_id`下二桁から特定した`partition_seed`をパーティションキーとして使用した投稿用テーブル。

特定ユーザーの投稿に特化しています。

### 主キー
- `post_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| user_id     | UUID      | NOT NULL, PRIMARY KEY | ユーザーID |
| post_id     | BIGINT    | NOT NULL, PRIMARY KEY  | 投稿ID、SnowflakeIDを使用 |
| post_type | TEXT | NOT NULL | 投稿タイプ(fuel,eventなど) |
| vehicle_id | TEXT |  | 車両ID、NULL許容 |
| reply_to | BIGINT |  | 返信先のpost_id |
| repost_to | BIGINT |  | リポスト先のpost_id |
| like_num | INT |  | いいねされている数 |
| reply_num | INT |  | リプライされている数 |
| repost_num | INT |  | リポストされている数 |
| text | TEXT | NOT NULL | 投稿本文 |
| image_folder_url | TEXT |  | 画像フォルダのURL |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、user_idのUUID右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

投稿ID(post_id)といいねした日時(liked_at)に基づくインデックスを設定して、効率的なデータ取得を実現します。

- `CREATE INDEX idx_posts_on_user_id_partition_seed_post_id ON posts (user_id, partition_seed, post_id DESC);`

### 制約と推奨事項

- `post_id`を使ってユニークにしている。
- `user_id`がuser_basics、user_detailsテーブルの主キー
- `post_id`はCassandra内のfuels、eventsテーブルの外部キー

### クエリパターン

特定のユーザーの最新の投稿を取得するクエリ：

```sql
SELECT * FROM posts
WHERE user_id = '特定のuser_id'
AND partition_seed = 特定のpartition_seed
ORDER BY post_id DESC
LIMIT 10;
```

### テーブル定義

```sql
CREATE TABLE posts (
    user_id UUID NOT NULL,
    post_id BIGINT NOT NULL,
    post_type TEXT NOT NULL,
    vehicle_id TEXT,
    reply_to BIGINT,
    repost_to BIGINT,
    like_num INT NOT NULL,
    reply_num INT NOT NULL,
    repost_num INT NOT NULL,
    text TEXT NOT NULL,
    image_folder_url TEXT,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (post_id, partition_seed)
    ) PARTITION BY LIST (partition_seed);
```

## テーブル名: `likes`

### 概要
`post_id`に基づくいいね機能用テーブル。
リストパーティショニングを用いて、Snowflake IDのシーケンス番号の右から8桁を基に256通りに分割する。

### 主キー
- `user_id`, `post_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| post_id     | BIGINT    | NOT NULL, PRIMARY KEY  | 投稿ID、SnowflakeIDを使用 |
| user_id     | UUID      | NOT NULL, PRIMARY KEY | いいねしたユーザーのUUID |
| liked_at    | TIMESTAMP WITH TIME ZONE | NOT NULL | いいねの日時 |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、post_idのSnowflake IDのシーケンス番号の右から8桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

投稿ID(post_id)といいねした日時(liked_at)に基づくインデックスを設定して、効率的なデータ取得を実現します。

- `CREATE INDEX idx_likes_on_partition_seed_post_id_liked_at ON likes (partition_seed, post_id, liked_at DESC);`

### 制約と推奨事項

- `user_id`と`post_id`を使ってユニークにしている。
- `user_id`がuser_basics、user_detailsテーブルの主キー
- `post_id`はCassandra内のposts、posts_monthテーブルの主キー

### クエリパターン

特定の投稿の最新のいいねを取得するクエリ：

```sql
SELECT * 
FROM likes 
WHERE partition_seed = '計算された値' AND post_id = '特定のpost_id' 
ORDER BY liked_at DESC;
```

### テーブル定義

```sql
CREATE TABLE likes (
    post_id BIGINT NOT NULL,
    user_id UUID NOT NULL,
    liked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (post_id, user_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `followings`

### 概要
特定のユーザー(user_id)によるフォロー管理用テーブル。
リストパーティショニングを用いて、UUIDの下二桁(16進数)を基に256通りのパーティションに分割保存する。

### 主キー
- `user_id`, `following_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| user_id     | UUID      | NOT NULL, PRIMARY KEY | フォローしたユーザーのUUID |
| following_id    | UUID | NOT NULL, PRIMARY KEY | フォローされたユーザーのUUID            |
| following_at    | TIMESTAMP WITH TIME ZONE | NOT NULL  | フォローした日時            |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、user_id(uuid)の右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

フォローの日時(followed_at)に基づくインデックスを設定して、効率的なデータ取得を実現します。

- `CREATE INDEX idx_followings_on_partition_seed_user_id_followed_at ON followings (partition_seed, user_id, followed_at DESC);`

### 制約と推奨事項

- `user_id`と`following_id`を使ってユニークにしている。
- `user_id`がuser_basics、user_detailsテーブルの主キー
- 下記followersテーブルと一貫性を保つ必要あり

### クエリパターン

特定のユーザーの最新の投稿を取得するクエリ：

```sql
SELECT following_id
FROM followings
WHERE partition_seed = '計算された値' user_id = '特定のuser_id'
ORDER BY followed_at DESC;
```

### テーブル定義

```sql
CREATE TABLE followings (
    user_id UUID NOT NULL,
    following_id TEXT NOT NULL,
    followed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (user_id, following_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `followers`

### 概要
特定のユーザー(user_id)によるフォロ「ワ」ー管理用テーブル。
リストパーティショニングを用いて、UUIDの下二桁(16進数)を基に256通りのパーティションに分割保存する。

### 主キー
- `user_id`, `follower_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| user_id     | UUID      | NOT NULL, PRIMARY KEY | フォローされたユーザーのUUID |
| follower_id | UUID | NOT NULL, PRIMARY KEY | 急にフォローしてきたユーザーのUUID            |
| followed_at    | TIMESTAMP WITH TIME ZONE | NOT NULL  | フォローされた日時            |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、user_id(uuid)の右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

フォローの日時(followed_at)に基づくインデックスを設定して、効率的なデータ取得を実現します。

- `CREATE INDEX idx_followers_on_partition_seed_user_id_followed_at ON followers (partition_seed, user_id, followed_at DESC);`

### 制約と推奨事項

- `user_id`と`follower_id`を使ってユニークにしている。
- `user_id`がuser_basics、user_detailsテーブルの主キー
- 上記followingテーブルと一貫性を保つ必要あり

### クエリパターン

特定のユーザーの最新の投稿を取得するクエリ：

```sql
SELECT follower_id
FROM followers
WHERE partition_seed = '計算された値' user_id = '特定のuser_id'
ORDER BY followed_at DESC;
```

### テーブル定義

```sql
CREATE TABLE followers (
    user_id UUID NOT NULL,
    follower_id TEXT NOT NULL,
    followed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (user_id, follower_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `user_basics`

### 概要
user_idに紐づけられた、基本的な情報を格納。
投稿の表示時に最低限必要なものをカラムにします。
user_id下二桁(16進数で256通り)のパーティションシードによって分散。

### 主キー
- `user_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| user_id     | UUID      | NOT NULL, PRIMARY KEY | ユーザー固有のID |
| username    | TEXT | NOT NULL  | アルファベット数字のみ |
| handlename    | TEXT | NOT NULL  | ハンネ |
| user_images_url | TEXT | NOT NULL  | AWSのユーザー画像用フォルダ |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、user_id(uuid)の右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

`user_id`のみに基づくインデックスを設定。

- 今のところなし

### 制約と推奨事項

- `user_id`をUUIDとし、ユニークにしています。
- `user_id`がuser_detailsテーブルの主キー。
- user_detailsテーブル使用時に必ず使用される。

### クエリパターン

特定のユーザーの基本情報を取得：

```sql
SELECT user_id, username, handlename, user_images_url
FROM user_basics
WHERE partition_seed = '計算または指定された値' AND user_id = '指定のUUID';
```

複数のユーザーの基本情報を一度に取得：

```sql
SELECT user_id, username, handlename, user_images_url
FROM user_basics
WHERE partition_seed = '計算または指定された値' AND user_id IN ('user_id1', 'user_id1', 'user_id1'...);
```

### テーブル定義

```sql
CREATE TABLE user_basics (
    user_id UUID NOT NULL,
    username TEXT NOT NULL,
    handlename TEXT NOT NULL,
    user_images_url TEXT NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (user_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `user_details`

### 概要
user_idに紐づけられた、詳細な情報を格納。
主にプロフィール画面の表示時に最低限必要なものをカラムにします。
user_id下二桁(16進数で256通り)のパーティションシードによって分散。

### 主キー
- `user_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| user_id     | UUID      | NOT NULL, PRIMARY KEY | ユーザー固有のID |
| joined_at | TIMESTAMP WITH TIME ZONE | NOT NULL  | ユーザーがアカウントを作成した日時 |
| following_num  | INT | NOT NULL  | ユーザーがフォローしている人数 |
| follower_num  | INT | NOT NULL  | ユーザーのフォロワーの人数 |
| location | TEXT |  | ユーザーの位置情報 |
| simulator | TEXT |  | 最も使っているシミュレーター |
| user_bio_text | TEXT |  | ユーザーの自己紹介文 |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、user_id(uuid)の右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

`user_id`のみに基づくインデックスを設定。

- 今のところなし

### 制約と推奨事項

- `user_id`をUUIDとし、ユニークにしています。
- `user_id`がuser_basicsテーブルの主キー
- 使用時、必ずuser_basicテーブルが必要

### クエリパターン

特定のユーザーの詳細情報を取得：

```sql
SELECT user_id, joined_at, following_num, follower_num, location, has_simulator, user_bio_text
FROM user_details
WHERE partition_seed = '計算または指定された値' AND user_id = '指定されたUUID';
```

### テーブル定義

```sql
CREATE TABLE user_details (
    user_id UUID NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL,
    following_num INT NOT NULL,
    follower_num INT NOT NULL,
    location TEXT,
    simulator TEXT,
    user_bio_text TEXT,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (user_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `vehicle_basics`

### 概要
vehicle_idに紐づけられた、基本的な情報を格納。
投稿の表示時に最低限必要なものをカラムにします。
vehicle_id下二桁(16進数で256通り)のパーティションシードによって分散。

### 主キー
- `vehicle_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| vehicle_id     | UUID      | NOT NULL, PRIMARY KEY | 車両固有のID |
| manufacture | TEXT |  | メーカー名 |
| model_type  | TEXT |  | 型式 |
| nickname    | TEXT |  | 愛称 |
| vehicle_images_url | TEXT | NOT NULL  | AWSの車両基本画像用フォルダ |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、vehicle_id(uuid)の右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

`vehicle_id`のみに基づくインデックスを設定。

- `CREATE INDEX idx_model_type ON vehicles (model_type);`

### 制約と推奨事項

- `vehicle_id`をUUIDとし、ユニークにしています。
- `vehicle_id`がvehicle_detailsテーブルの主キー
- vehicle_detailsテーブル使用時に必ず使用される。

### クエリパターン

特定の車両の基本情報を取得：

```sql
SELECT vehicle_id, model_type, nickname, vehicle_images_url
FROM vehicle_basics
WHERE vehicle_id = '指定されたUUID';
```

モデルタイプに基づく車両情報の取得：

```sql
SELECT vehicle_id, model_type, nickname, vehicle_images_url
FROM vehicle_basics
WHERE model_type = '指定された型式';
```

### テーブル定義

```sql
CREATE TABLE vehicle_basics (
    vehicle_id UUID NOT NULL,
    manufacture TEXT,
    model_type TEXT,
    nickname TEXT,
    vehicle_images_url TEXT NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (vehicle_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `vehicle_details`

### 概要
vehicle_idに紐づけられた、詳細の情報を格納。
乗り物詳細表示時などに必要な、いろんな情報をカラムにします。
vehicle_id下二桁(16進数で256通り)のパーティションシードによって分散。

### 主キー
- `vehicle_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| vehicle_id     | UUID      | NOT NULL, PRIMARY KEY | 車両固有のID |
| joined_at | TIMESTAMP WITH TIME ZONE | NOT NULL | project2登録日 |
| manufacture_date | TIMESTAMP WITH TIME ZONE |  | 製造年月（日はない） |
| newcar_registration_date | TIMESTAMP WITH TIME ZONE | | 新車登録日 |
| status | TEXT | NOT NULL  | 稼働中、車検切れ、廃車など |
| mileage | INT |  | 総走行距離(km) |
| fan_num | TEXT | NOT NULL  | ファン（人）の数 |
| vehicle_bio_text | TEXT |  | 乗り物自由記述欄 |
| frame_no | TEXT |  | フレーム、及びシャシーナンバー |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、vehicle_id(uuid)の右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

`vehicle_id`のみに基づくインデックスを設定。

- 今のところなし

### 制約と推奨事項

- `vehicle_id`をUUIDとし、ユニークにしています。
- `vehicle_id`がvehicle_basicsテーブルの主キー
- 使用時、必ずvehicle_basicテーブルが必要

### クエリパターン

特定の車両の詳細情報を取得：

```sql
SELECT *
FROM vehicle_details
WHERE vehicle_id = '指定されたUUID';
```

### テーブル定義

```sql
CREATE TABLE vehicle_details (
    vehicle_id UUID NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL,
    manufacture_date TIMESTAMP WITH TIME ZONE,
    newcar_registration_date TIMESTAMP WITH TIME ZONE,
    status TEXT NOT NULL,
    mileage INT,
    fan_num INT NOT NULL,
    vehicle_bio_text TEXT,
    frame_no TEXT,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (vehicle_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `fuels`

### 概要
特定の車両に関する燃料記録を保持するテーブル。`vehicle_id` と vehicleID由来の `partition_seed` を使ってデータのパーティショニングを行い、特定の車両に対するリクエストを高速化します。

### 主キー
主キー構成: vehicle_id, partition_seed, refuel_at
### カラム
| Column Name     | Data Type | Constraints | Description               |
|-----------------|-----------|-------------|---------------------------|
| vehicle_id      | UUID      | NOT NULL    | 車両ID                    |
| fuel_id         | BIGINT    | NOT NULL    | 燃料記録のID (Snowflake)  |
| post_id         | BIGINT    | NOT NULL    | 関連投稿のID (Snowflake)  |
| refuel_at       | TIMESTAMP | NOT NULL    | 給油日時 (RFC3339)        |
| amount_ml       | INT       | NOT NULL    | 給油量 (単位はミリリットル) |
| fuel_type       | TEXT      | NOT NULL    | 燃料の種類 (例: diesel)   |
| liter_fee       | INT       |             | リッター当たりの価格       |
| total_fee       | INT       |             | 今回の給油にかかった総費用 |
| odometer_kilo   | INT       |             | 給油時の総走行距離         |
| tripmeter_kilo  | INT       |             | 最後の給油からの走行距離   |
| location        | TEXT      |             | 給油場所                  |
| receipt_img_url | TEXT      |             | レシート画像のURL         |
| partition_seed  | TEXT      | NOT NULL    | パーティション用バケット   |

### インデックス
```sql
CREATE INDEX idx_fuel_vehicle_time ON fuels (vehicle_id, partition_seed, refuel_at DESC);
```

### 制約
SQL文で WHERE 句のフィルタリングを使用する際には、インデックスを適切に使用し、パフォーマンスを維持するように設計します。

### クエリパターン
特定車両の最新の給油記録を取得するクエリ:
```sql
SELECT * FROM fuels
WHERE vehicle_id = ? AND partition_seed = ?
ORDER BY refuel_at DESC
LIMIT 10;
```

特定の期間の給油情報を取得するクエリ:
```sql
SELECT * FROM fuels
WHERE vehicle_id = ? AND partition_seed = ? AND refuel_at > '開始日時' AND refuel_at < '終了日時'
ORDER BY refuel_at DESC;
```

### テーブル定義
```sql
CREATE TABLE fuels (
    vehicle_id UUID NOT NULL,
    partition_seed TEXT NOT NULL,
    fuel_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    refuel_at TIMESTAMP NOT NULL,
    amount_ml INT NOT NULL,
    fuel_type TEXT NOT NULL,
    liter_fee INT,
    total_fee INT,
    odometer_kilo INT,
    tripmeter_kilo INT,
    location TEXT,
    receipt_img_url TEXT,
    PRIMARY KEY (vehicle_id, partition_seed, refuel_at)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `events`

### 概要
特定のユーザーID `event_id` とuser_idから算出した `partition_seed` に基づいてイベントを管理するためのテーブル。PostgreSQLではパーティションを活用してデータの管理を行うことができます。

### 主キー
- Partition Key: `event_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| user_id        | UUID       | NOT NULL               | ユーザーID              |
| partition_seed | TEXT       | NOT NULL               | パーティション用シード値 |
| event_id       | BIGINT     | NOT NULL               | イベントID (Snowflake)  |
| post_id        | UUID       | NOT NULL               | 関連投稿ID              |
| start_at       | TIMESTAMP  | NOT NULL               | イベント開始時間        |
| end_at         | TIMESTAMP  | NOT NULL               | イベント終了時間        |
| event_title    | TEXT       | NOT NULL               | イベントタイトル        |
| is_allday      | BOOLEAN    | NOT NULL               | 終日イベントかどうか    |
| location       | TEXT       |                        | イベント開催場所        |


```sql
CREATE TABLE events (
    user_id UUID NOT NULL,
    partition_seed TEXT NOT NULL,
    event_id BIGINT NOT NULL,
    post_id UUID NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    event_title TEXT NOT NULL,
    is_allday BOOLEAN NOT NULL,
    location TEXT,
    PRIMARY KEY (event_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `laptime_details`

### 概要
laptime_idに紐づけられた、基本的な情報を格納。
投稿の表示時に最低限必要なものをカラムにします。
laptime_idであるSnowflakeIDを256で余剰し、0～255の256通りのパーティションシードに分割

なお、basicsに関してはモデルの型式をユーザーの手入力ではなく、選んでもらうようにしてもらった後に最適化用に作成予定。
つまり、車両モデルのリファレンステーブルができた後だから当分まだかな…

### 主キー
- `laptime_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| laptime_id | BIGINT | NOT NULL, PRIMARY KEY | タイムのID(Snowflake) |
| track_id  | UUID | NOT NULL  | track_id。layoutと共用でレギュレーションを特定 |
| layout_id  | UUID | NOT NULL  | track_idとJOIN |
| user_id | UUID | NOT NULL  | ユーザーID |
| vehicle_id | UUID | NOT NULL  | このままだと、記録時の車両が残せない |
| laptime_ms | INT | NOT NULL  | ラップタイムはミリ秒で保存 |
| road_condition | UUID | NOT NULL  | dry, wet, half_wetとか |
| record_at | TIMESTAMP | NOT NULL  | RFC3339 |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値  |

### パーティショニング

このテーブルは、laptime_idであるSnowflakeIDを256で余剰し、0～255の256通りのパーティションシードに分割


### インデックス

`laptime_id`のみに基づくインデックスを設定。

- `CREATE INDEX idx_laptime_ms ON laptime_details(laptime_ms);`

### 制約と推奨事項

- `laptime_id`をSnowflake(uuid)とし、ユニークにしています。
- `user_id`がuser_basics、user_detailsテーブルの主キー
- `vehicle_id`がvehicle_basics、vehicle_detailsテーブルの主キー
- `layout_id`内の`track_id`とによって`tracks`テーブルとJOINする。

### クエリパターン

特定トラック、特定レイアウトでの最速ラップタイム:

```sql
SELECT MIN(laptime_ms) AS fastest_laptime
FROM laptime_details
WHERE layout_id = '特定のレイアウトUUID' AND partition_seed = "シード値"
```

### テーブル定義

```sql
CREATE TABLE laptime_details (
    laptime_id BIGINT NOT NULL,
    track_id UUID NOT NULL,
    layout_id UUID NOT NULL,
    user_id UUID NOT NULL,
    vehicle_id UUID NOT NULL,
    laptime_ms INT NOT NULL,
    road_condition UUID NOT NULL,
    record_at TIMESTAMP NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (laptime_id, partition_seed)
);
```

## テーブル名: `vehicle_user_links`

### 概要
`user_id`、`vehicle_id`による多対多の中間テーブルであり、ユーザーによる車両の所有関係を示す。

このテーブルにより、共同所有や車との関係、その関係の開始、終了時期、次の人への受け渡し方などを表現する。

本テーブルにおけるパーティショニングでは、vehicle_idをパーティションキーとして用いることで、歴代の所有者や共同所有なども含めたユーザーを効率的に処理する。

vehicle_id下二桁(16進数で256通り)のパーティションシードによって分散。

### 主キー
- `user_id`, `vehicle_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| user_id | UUID | NOT NULL, PRIMARY KEY | ユーザー固有のID |
| vehicle_id | UUID | NOT NULL, PRIMARY KEY | 車両固有のID |
| relation_type  | TEXT | NOT NULL | `own`(複数不可), `lease`, `shared`, `corporate`, `bye`など|
| start_at | TIMESTAMP WITH TIME ZONE | NOT NULL | 所有開始日時 |
| end_at | TIMESTAMP WITH TIME ZONE |  | 手放した日時 |
| price | INT | NOT NULL | お迎えしたときの金額 |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | レコード作成日 |
| access_level | TEXT | NOT NULL | `admin`, `moderator`など |
| transfer_type | TEXT | NOT NULL | `buy`, `gift`など |
| mileage_at_registration | BIGINT | NOT NULL | 登録時の走行距離 |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値 |

### パーティショニング

このテーブルは、vehicle_id(uuid)の右から2桁を使用してリストパーティショニングを行い、256通りにデータを分割します。

### インデックス

今のところなし

- ``

### 制約と推奨事項

- `user_id`, `vehicle_id`の組み合わせによりユニークにしています。
- このテーブル自体が`user_id`, `vehicle_id`を用いるそれぞれ２テーブルの中間テーブルです。
- `user_id`がuser_basics, user_detailsテーブルの主キー
- `vehicle_id`がvehicle_basics, vehicle_detailsテーブルの主キー

### クエリパターン

特定のユーザーが所有する車両を全て取得：
```sql
SELECT vehicle_id, relation_type, start_at, end_at
FROM vehicle_user_links
WHERE user_id = '特定のuser_id' AND end_at IS NULL;
```

特定の車両の歴代の所有者を取得：
```sql
SELECT user_id, relation_type, start_at, end_at
FROM vehicle_user_links
WHERE vehicle_id = '特定のvehicle_id'
ORDER BY start_at DESC;
```

### テーブル定義

```sql
CREATE TABLE vehicle_user_links (
    user_id UUID NOT NULL,
    vehicle_id UUID NOT NULL,
    relation_type TEXT NOT NULL,
    start_at TIMESTAMP WITH TIME ZONE NOT NULL,
    end_at TIMESTAMP WITH TIME ZONE,
    price INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    access_level TEXT NOT NULL,
    transfer_type TEXT NOT NULL,
    mileage_at_registration BIGINT NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (user_id, vehicle_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `tracks`

### 概要
`track_id`による、サーキットの大まかな情報を示すテーブル。

後述の`layouts`テーブルと使用することで、実質的に継承のような役割を果たす。

track_id下一桁(16進数で16通り)のパーティションシードによって分散。

### 主キー
- `track_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| track_id | UUID | NOT NULL, PRIMARY KEY | トラック固有のID |
| user_id | UUID | NOT NULL | ユーザーID |
| track_name | TEXT | NOT NULL | サーキットの名前 |
| country | TEXT | NOT NULL | 所在する国名 |
| location | TEXT | NOT NULL | 所在する場所の名前。柔軟に。 |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値 |

### パーティショニング

このテーブルは、track_id(uuid)の右から1桁を使用してリストパーティショニングを行い、16通りにデータを分割します。

### インデックス

今のところなし

- ``

### 制約と推奨事項

- 

### クエリパターン

特定の`track_id`に対応するコース情報を全て取得：
```sql
    SELECT *
    FROM tracks
    WHERE track_id = '指定のtrack_id値'
    AND partition_seed = '指定のpartition_seed値';
```

### テーブル定義

```sql
CREATE TABLE tracks (
    track_id UUID NOT NULL,
    user_id UUID NOT NULL,
    track_name TEXT NOT NULL,
    country TEXT NOT NULL,
    location TEXT NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (track_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```

## テーブル名: `layouts`

### 概要
`layout_id`による、サーキットの特定のレイアウト情報を示すテーブル。

`track_id`下一桁(16進数で16通り)のパーティションシードによって分散。

### 主キー
- `track_id`, `partition_seed`

### カラム

| Column Name | Data Type  | Constraints            | Description            |
|-------------|------------|------------------------|------------------------|
| layout_id | UUID | NOT NULL, PRIMARY KEY | レイアウト固有のID |
| track_id | UUID | NOT NULL | トラック固有のID |
| layout_name | TEXT | NOT NULL| レイアウトの名前 |
| latitude | DECIMAL(9,6) | NOT NULL | 緯度。経度と合わせてGoogleMapsAPI用のジオタグとする |
| longitude | DECIMAL(9,6) | NOT NULL | 経度 |
| length | INT | NOT NULL | そのレイアウトでの全長。メートルで保存 |
| partition_seed | TEXT | NOT NULL, PRIMARY KEY | パーティショニング用のシード値 |

### パーティショニング

このテーブルは、track_id(uuid)の右から1桁を使用してリストパーティショニングを行い、16通りにデータを分割します。

### インデックス

今のところなし

- ``

### 制約と推奨事項

- 

### クエリパターン

特定の`track_id`、`layout_id`に対するレイアウトの詳細：
```sql
    SELECT *
    FROM layouts
    WHERE track_id = '指定したトラックのUUID'
    AND partition_seed = '特定のパーティションシード';
```

### テーブル定義

```sql
CREATE TABLE layouts (
    layout_id UUID NOT NULL,
    track_id UUID NOT NULL,
    layout_name TEXT NOT NULL,
    latitude DECIMAL(9,6) NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    length INT NOT NULL,
    partition_seed TEXT NOT NULL,
    PRIMARY KEY (layout_id, partition_seed)
) PARTITION BY LIST (partition_seed);
```