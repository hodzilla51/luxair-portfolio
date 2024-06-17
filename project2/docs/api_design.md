FORMAT: 1A
HOST: http://

# API設計書

## 概要

ユーザー情報の作成、変更、削除用のエンドポイント。
Keycloakとの連携が必要です。

| ホスト | プロトコル |バージョン| データ形式 |
|-----------|-----------|-----------|-----------|
| localhost:5000  | http |v1| JSON |

http://localhost:5000/api/v1

顔のマークがついているところはとりあえず実装されているところです🥺

## アクセストークン

必要なAPIリクエストのAuthorizationリクエストヘッダに、Keycloakのアクセストークンを付与してください。

```例    
Authorization: Bearer secretaccesstoken
```

## ステータスコード

| ステータスコード | 説明                       |
|-----------------|----------------------------|
| 200             | OK - リクエストが成功したことを示します。 |
| 201             | Created - リソースが正常に作成されたことを示します。 |
| 400             | Bad Request - サーバーがリクエストを理解できなかったことを示します。 |
| 401             | Unauthorized - 認証が必要なことを示します。 |
| 403             | Forbidden - サーバーがリクエストを拒否したことを示します。 |
| 404             | Not Found - リクエストしたリソースが見つからなかったことを示します。 |
| 429             | Too Many Requests - リクエストが多すぎることを示します。 |
| 500             | Internal Server Error - サーバー内部でエラーが発生したことを示します。 |
| 503             | Service Unavailable - サーバーがリクエストを処理できない状態であることを示します。 |

## JSONの取り扱いについて（2024.3.19追記）
JSONは全てキャメルケースにて統一します。
途中までテキトーだったので、スネークケースなどが混在している恐れがありますので、発見し次第、キャメルケースに変更してください。

- 正：`camelCase`
- 誤：`snake_case`

## エンドポイント [/login]

### 登録済みユーザーのログイン [POST] [/login]🥺😵‍💫

POSTされたusername、email、パスワードを元に、Keycloakに認証トークンを問い合わせます。

+ Request (application/json)
```json
    {
        "username": "HelloDog",
        "email": ,
        "password": 
    }
```

+ Response 200 OK (application/json)
```json
    {
        "accessToken": "取得したアクセストークン"
    }
```
+ Response 401 Unauthorized (application/json)
```json
    {
        "error": "認証に失敗しました。ユーザー名またはパスワードが正しくありません。"
    }
```

状態管理：localStorage, zustand



## エンドポイント [/users]

### ユーザー情報の新規登録 [POST] [/users]🥺😵‍💫

ユーザーを新規作成します。誰でもできます。
パスワードは`^[a-zA-Z0-9!@#$%^&*()_+{}":;'<>?,.\/-]+$`でバリデーション。usernameは`^[a-zA-Z0-9]+$`。

現在、実装の都合上エンドポイントが複雑化しています。下のjsonも形式が変わっています。あくまで参考程度にして、詳細はコードを参照してください。(2024.4.3追記)

+ Request (application/json)
```json
    {
        "username": "HelloDog",
        "handlename": "こんにちは犬",
        "email": "hello@dog.com",
        "password": "password123",
        "phoneNumber": "123-456-7890",
        "simulator": "Assetto corsa",
        "userBioText": "こんにちは、私はこんにちは犬です！",
        "userIconJPG": "なんかbase64の長いやつ",
        "userHeaderJPG": "なんかbase64の長いやつ"
    }
```

+ Response 201 (application/json)
```json
        {
            "message": "ユーザーが正常に作成されました。",
        }
```

状態管理：useState

### フォロー関係の新規登録（トグル） [POST] [/users/follow]🥺

フォロー関係をトグルで`POST`、`DELETE`します。
ユーザー毎のプロフィール画面にフォローボタンを設置してください。

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)

```json
        {
            "followingId": "対象のユーザーID"
        }
```

フォローした場合のレスポンス：

+ Response 201 (application/json)

```json

        {
            "message": "フォローしました"
        }
```

フォローを解除した場合のレスポンス：
```json
        {
            "message": "フォローを解除しました"
        }
```

状態管理：useState

### ユーザーのブロック関係の新規登録（トグル） [POST] [/users/block]まだいい

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)

        {
            "暇なとき書いてね"
        }


+ Response 201 (application/json)

        {
            "message": "ユーザーが正常に作成されました。",
        }

### ユーザーの基本情報の取得 [GET] [/users/{userID}/basic]🥺😵‍💫

ユーザーIDに対応するユーザーの基本情報を取得します。
これは通常、投稿データなどに付随するので必要ないですが、何らかのエラーでプロフィール画面などで基本情報を失った際に再度要求するためのエンドポイントです。

+ Parameters
    + userID (required, string) - 取得したいユーザーのID

+ Response 200 (application/json)
```json
        {
                "username": "ESPDog",
                "handlename": "天才",
                "userImagesURL": "なんかbase64の長いやつ"
        }
```

状態管理：SWR

### ユーザーの詳細情報の取得 [GET] [/users/{userID}/detail]🥺😵‍💫

ユーザーIDに対応するユーザーの詳細情報を取得します。
これはプロフィールページの表示に使用されます。

+ Parameters
    + userId (required, string) - 取得したいユーザーのID

+ Response 200 (application/json)
```json
        {
            "joinedAt": "2023/11/25",
            "followingNum": 13,
            "followerNum": 15,
            "location": "東京",
            "simulator": "Assetto corsa",
            "userBioText": "どうも自己紹介文です。50文字くらいまで書けるようにしたいです。",
        }
```

状態管理：SWR

### ユーザーの所有する車両一覧を取得 [GET] [/users/{userID}/vehicles]🥺

`userID`に対応する車両の一覧を取得します。
これはプロフィールページの車両カード表示に使用されます。

～バックエンド向け～
`userID`を受け取った後、`vehicle_user_links`と`vehicle_basics`の内部結合で、
`vehicle_id`, `model_type`, `nickname`, `vehicle_images_url`を配列として取得し、レスポンスしてください。
例：
```sql
    SELECT
        v.vehicle_id AS "vehicleID",
        v.model_type,
        v.nickname,
        v.vehicle_images_url,
        vur.relation_type
    FROM
        vehicles v
    JOIN
        vehicle_user_links vur ON v.vehicle_id = vur.vehicle_id
    WHERE
        vur.user_id = '指定されたuserID'
```

+ Parameters
    + userId (required, string) - 取得したいユーザーのID

+ Response 200 (application/json)
```json
    {
        "userID": "",
        "vehicles": [
            {
                "vehicleID": "uuid-1",
                "modelType": "JZX-100",
                "nickname": "チェイサー",
                "vehicleImagesURL": "https://～（AWS）"
            },
            {
                "vehicleID": "uuid-2",
                "modelType": "LP-640",
                "nickname": "ウラカン",
                "vehicleImagesURL": "https://～（AWS）"
            },
            {
                "vehicleID": "uuid-3",
                "modelType": "",
                "nickname": "ヨタハチ",
                "vehicleImagesURL": "https://～（AWS）"
            }
        ]
    }
```
状態管理：SWR

### ユーザーの予定一覧を取得 [GET] [/users/{userID}/calendar]できたら

### ユーザーのラップタイム一覧を取得 [GET] [/users/{userID}/laptimes]🥺

`userID`に対応するラップタイムの一覧を取得します。
これはプロフィールページのラップタイム欄の表示（仮）に使用されます。

+ Parameters
    + userId (required, string) - 取得したいユーザーのID

+ Response 200 (application/json)
```json
    {
        "userID": "",
        "laptimes": [
            {
                "laptimeID": 123456789012345678,
                "trackID": "123e4567-e89b-12d3-a456-426614174000",
                "layoutID": "123e4567-e89b-12d3-a456-426614174000",
                "userID": "123e4567-e89b-12d3-a456-426614174000",
                "vehicleID": "123e4567-e89b-12d3-a456-426614174000",
                "laptimeMs": 85000,
                "roadCondition": "123e4567-e89b-12d3-a456-426614174000",
                "recordAt": "2023-03-11T14:20:00Z"
            },
            {
                "laptimeID": 123456789012345678,
                "trackID": "123e4567-e89b-12d3-a456-426614174000",
                "layoutID": "123e4567-e89b-12d3-a456-426614174000",
                "userID": "123e4567-e89b-12d3-a456-426614174000",
                "vehicleID": "123e4567-e89b-12d3-a456-426614174000",
                "laptimeMs": 85000,
                "roadCondition": "123e4567-e89b-12d3-a456-426614174000",
                "recordAt": "2023-03-11T14:20:00Z"
            },
            {
                "laptimeID": 123456789012345678,
                "trackID": "123e4567-e89b-12d3-a456-426614174000",
                "layoutID": "123e4567-e89b-12d3-a456-426614174000",
                "userID": "123e4567-e89b-12d3-a456-426614174000",
                "vehicleID": "123e4567-e89b-12d3-a456-426614174000",
                "laptimeMs": 85000,
                "roadCondition": "123e4567-e89b-12d3-a456-426614174000",
                "recordAt": "2023-03-11T14:20:00Z"
            }
        ]
    }
```

状態管理：SWR

### ユーザーのいいねした投稿の一覧を取得 [GET] [/users/{userID}/likes]一旦保留


`userID`に対応するいいねの一覧を取得します。
これはプロフィールページのいいね欄の表示（仮）に使用されます。

+ Parameters
    + userId (required, string) - 取得したいユーザーのID

+ Response 200 (application/json)
```json
    [
        {
            "likes": [
            {
                "userID": "post123", // 実際はuuid
                "likedAt": "RFC3339",
            },
            {
                "userID": "post123", // 実際はuuid
                "likedAt": "RFC3339",
            }
            // この配列はlimitパラメータによって返される投稿の数が増減します
            ],
            "NewLastLikedAt": ""
        }
    ]
```

### ユーザーのフォロー関係の一覧を取得 [GET]　[/users/{userID}/followings]🥺

ユーザーIDに対応するフォロー関係を取得します。
これはプロフィールページのフォロー欄に使用されます。

+ Parameters
    + user_id (required, string) - 取得したいユーザーのID。
    + last_following_at (required, string) - 取得した最後のフォローの時間。
    + limit (optional, string) - 取得する数。デフォルトは30。


+ Response 200 (application/json)

```json
    {
        "userID": "",
        "followings": [
            {
                "userID": "user123",
                "followingID": "user456",
                "followedAt": "2023-11-25T15:30:00Z"
            },
            {
                "userID": "user123",
                "followingID": "user789",
                "followedAt": "2023-11-26T16:45:00Z"
            },
            {
                "userID": "user123",
                "followingID": "user101",
                "followedAt": "2023-11-27T11:20:00Z"
            }
            // この配列はlimitパラメータによって返される投稿の数が増減します
        ]
    }
```

+ Response 429 (application/json)
```json
        {
            "error": "Too Many Requests",
            "message": "リクエスト数が多すぎます。しばらく待ってから再度試してください。"
        }
```
状態管理：SWR

### ユーザーのフォロワー関係の一覧を取得 [GET]　[/users/{userID}/followers]🥺

ユーザーIDに対応するフォロワー関係を取得します。
これはプロフィールページのフォロワー欄に使用されます。

+ Parameters
    + user_id (required, string) - 取得したいユーザーのID
    + limit (optional, string) - 取得する数。デフォルトは30。

+ Response 200 (application/json)
```json
        [
            {
                "userID": "user456",
                "followerID": "user123",
                "followedAt": "2023-11-25T15:30:00Z"
            },
            {
                "userID": "user789",
                "followerID": "user123",
                "followedAt": "2023-11-26T16:45:00Z"
            },
            {
                "userID": "user101",
                "followerID": "user123",
                "followedAt": "2023-11-27T11:20:00Z"
            }
        ]
```

+ Response 201 (application/json)
```json
        {
            "message": "ユーザー情報の取得に成功。",
        }
```

状態管理：SWR

### ユーザー情報の更新 [PATCH] [/users/me]🥺

Authorizationされているユーザーのユーザー情報の更新をします。

～バックエンドの人へ～
`username`の更新だけは別個作ろうかと思っています。なので今回は無しで。
`handlename`とそれ以外は、それぞれbasicsとdetailsで分けられているので、結合なりをしてください。

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)
```json
        {
            "handlename": "新しいハンドル名",
            "location": "新しい位置情報",
            "userBioText": "新しい自己紹介文",
            "userIconURL": "",
            "userHeaderURL": ""
        }
```
+ Response 200 (application/json)
```json
    {
         "message": "ユーザー情報を更新しました。"
    }
```

状態管理：useState

### ユーザーのプライバシー設定の更新 [PATCH] [/users/me/privacy]まだいい

### ユーザーの通知設定の更新 [PATCH] [/users/me/notification]まだいい

### ユーザーの退会 [DELETE] [/users/me]まだいい

## エンドポイント [/vehicles]

### 車両情報の新規登録 [POST] [/vehicles]🥺

    車両の新規登録をします。
    リクエストには`vehicle_basics`, `vehicle_details`, `vehicle_user_links`の三テーブルが必要です。これらを結合するなどして全てに送ってください。

    なお、`mileage`の-1を走行距離不明、疑義車として扱います。
    `manufacture_date`、`newcar_registration_date`の場合、時間のゼロ値の場合が新車。nullを不明車とします。


+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)
```json
        {
                "modelType": "型式",
                "nickname": "愛称",
                "vehicleIconJPG": "http://aws.com/users/265c66a4-3f0a-4118-9a4f-56725d2c717b",
                "vehicleHeaderJPG": "http://aws.com/users/265c66a4-3f0a-4118-9a4f-56725d2c717b",
                "manufactureDate": "2023-05-01T00:00:00Z",
                "newcarRegistrationDate": "2023-06-15T00:00:00Z",
                "status": "active",
                "mileage": 1029,
                "fanNum": 0,
                "vehicleBioText": "乗り物自由記述欄",
                "frameNo": "フレームナンバーとか",
                "userID": "265c66a4-3f0a-4118-9a4f-56725d2c717b",
                "relationType": "own",
                "startAt": "2024-03-12T14:30:00Z",
                "price": 120387,
                "accessLevel": "admin",
                "transferType": "buy",
                "mileageAtRegistration": 1029
        }
```

+ Response 201 (application/json)
```json
        {
            "message": "車両が正常に登録されました。",
        }
```

状態管理：useState

### 車両の基本情報の取得 [GET] [/vehicles/{vehicleID}/basic]🥺

ユーザーIDに対応するユーザーの基本情報を取得します。
これは通常、ユーザーの車両一覧を取得する際などに付随するので必要ないですが、何らかのエラーで基本情報を失った際に再度要求するためのエンドポイントです。

+ Parameters
    + vehicleId (required, string) - 取得したい車両のID

+ Response 200 (application/json)
```json
    {
            "manufacture": "メーカー名",
            "modelType": "型式",
            "nickname": "愛称",
            "vehicleImagesURL": "AWSの車両基本画像用フォルダURL"
    }
```

状態管理：SWR

### 車両の詳細情報の取得 [GET] [/vehicles/{vehicleID}/detail]🥺

車両IDに対応する車両の詳細情報を取得します。
これは車両詳細表示(仮)の表示に使用されます。

+ Parameters
    + vehicleId (required, string) - 取得したい車両のID

+ Response 200 (application/json)
```json
    {
            "joinedAt": "project2登録日(例: 2024-03-12T14:30:00Z)",
            "manufactureDate": "製造年月(例: 2023-05-01T00:00:00Z)",
            "newcarRegistrationDate": "新車登録日(例: 2023-06-15T00:00:00Z)",
            "status": "稼働中",
            "mileage": "総走行距離(km)",
            "fanNum": "ファン（人）の数",
            "vehicleBioText": "乗り物自由記述欄",
            "frameNo": "フレーム、及びシャシーナンバー"
    }
```

状態管理：SWR

### 車両の関係ユーザーの一覧を取得 [GET] [/vehicles/{vehicleID}/users]🥺

車両IDに対応する車両と何らかの関係を持つユーザーを取得します。
これは車両詳細表示(仮)の表示に使用されます。

+ Parameters
    + vehicleId (required, string) - 取得したい車両のID

+ Response 200 (application/json)
```json
        {
        "vehicleID": "指定された車両ID",
        "users": [
            {
                "userID": "uuid-1",
                "username": "fafafa",
                "handlename": "ハンネ",
                "userImagesURL": "http",
                "relationType": "own",
                "startAt": "所有開始日時(例: 2024-03-12T14:30:00Z)",
                "endAt": "手放した日時(例: 2024-06-12T14:30:00Z)",
                "createdAt": "レコード作成日(例: 2024-03-01T10:00:00Z)",
                "accessLevel": "admin",
                "transferType": "buy",
                "mileageAt": "登録時の走行距離(例: '5000km')"
            },
            {
                "userID": "uuid-1",
                "username": "fafafa",
                "handlename": "ハンネ",
                "userImagesURL": "http",
                "relationType": "own",
                "startAt": "所有開始日時(例: 2024-03-12T14:30:00Z)",
                "endAt": "手放した日時(例: 2024-06-12T14:30:00Z)",
                "createdAt": "レコード作成日(例: 2024-03-01T10:00:00Z)",
                "accessLevel": "admin",
                "transferType": "buy",
                "mileageAt": "登録時の走行距離(例: '5000km')"
            },
        ]
    }
```

状態管理：SWR

### 車両の給油フィードを取得 [GET] [/vehicles/{vehicleID}/fuels]🥺

車両IDに対応する給油情報一覧を取得します。
ガソリン記録欄に使われます。ガソリン情報はsnowflakeです。

～バックエンド向け～
Cassandraを使います。十分に注意してください。

+ Parameters
    + vehicleId (required, string) - 取得したい車両のID

+ Response 200 (application/json)

```json
{
  "vehicleID": "指定された車両ID",
  "refuelingFeed": [
    {
      "fuelID": "snowflake",
      "postID": 1234567890123456,
      "refuelAt": "2024-03-12T14:30:00Z",
      "amountMl": 50000,
      "fuelType": "diesel",
      "literFee": 120,
      "totalFee": 6000,
      "odometerKilo": 105000,
      "tripmeterKilo": 500,
      "location": "東京",
      "receiptImgURL": "https://aws.example.com/receipt1.jpg"
    },
    {
      "fuelID": "uuid-2",
      "postID": 1234567890123457,
      "refuelAt": "2024-03-25T15:00:00Z",
      "amountMl": 45000,
      "fuelType": "diesel",
      "literFee": 115,
      "totalFee": 5175,
      "odometerKilo": 105500,
      "tripmeterKilo": 500,
      "location": "神奈川",
      "receiptImgURL": "https://aws.example.com/receipt2.jpg"
    }
  ]
}
```

状態管理：SWR

### 車両のラップタイム一覧を取得 [GET] [/vehicles/{vehicleID}/laptimes]🥺

車両IDに対応するラップタイム一覧を取得します。ラップタイム情報はsnowflakeです

+ Parameters
    + vehicleId (required, string) - 取得したい車両のID

+ Response 200 (application/json)

```json
{
  "vehicleID": "指定された車両ID",
  "lapTimes": [
    {
      "laptimeID": 1234567890123456,
      "trackID": "uuid-1",
      "layoutNum": 1,
      "userID": "uuid-2",
      "vehicleID": "指定された車両ID",
      "laptimeMs": 120000,
      "roadCondition": "dry",
      "recordAt": "2024-03-12T14:30:00Z"
    },
    {
      "laptimeID": 1234567890123457,
      "trackID": "uuid-3",
      "layoutNum": 2,
      "userID": "uuid-4",
      "vehicleID": "指定された車両ID",
      "laptimeMs": 118500,
      "roadCondition": "wet",
      "recordAt": "2024-03-25T15:00:00Z"
    }
  ]
}
```

状態管理：SWR


### 車両情報の更新 [PATCH] [/vehicles?vehicleID={vehicleID}]

Authorizationされているユーザーかつ、そのユーザーが車両の`accsess_level` `admin`を所持している車両の情報を更新します。

～バックエンドの人へ～
非常に重要な機能なので、セキュリティホールがないように作ってください。データベースへのアクセスは二度になることが推測されます。
以下の要件を満たしてください。

- Authorizationチェック:
リクエストを送るユーザーが認証されていること。
認証されたユーザーが指定された車両に対してadminレベルのアクセス権限を持っていることを確認する。

- 車両情報の更新権限チェック:
vehicle_user_linksテーブルを参照して、認証されたユーザーがaccess_levelがadminである車両の情報のみ更新できるようにする。

- データの更新:
ユーザーが権限を持っていることが確認できたら、vehicle_basicsおよびvehicle_detailsテーブルに対して、リクエストボディに含まれる更新データを適用する。
一部の更新が失敗した場合に全体の一貫性を保つため、更新操作はトランザクション管理下で行う。

- レスポンスの送信:
更新が成功した場合は、更新後の車両情報を含むレスポンスをクライアントに返す。
更新に必要な権限がない場合やその他のエラーが発生した場合は、適切なエラーメッセージとHTTPステータスコードを返す。

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)
```json
    {
        "vehicleID": "指定された車両ID",
        "updates": {
            "vehicleBasics": {
            "manufacture": "HONDA",
            "modelType": "新しい型式",
            "nickname": "新しい愛称",
            "vehicleIconJpg": "",
            "vehicleHeaderJpg": ""
            },
            "vehicleDetails": {
            "status": "稼働中",
            "mileage": 150000,
            "vehicleBioText": "新しい乗り物自由記述欄",
            "frameNo": "新しいフレーム、及びシャシーナンバー"
            }
        }
    }
```
+ Response 200 (application/json)
```json
     "message": "ユーザー情報を更新しました。"
```

状態管理：useState

## エンドポイント [/posts]

### 投稿情報の新規作成 [POST] [/posts]🥺

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json) 化け物か？
```json
        {
            "text": "投稿テキスト",
            "vehicleID": "車両ID",
            "repostBy": 0,
            "replyTo": 0,
            "postType": "ポストタイプ",
            "imageJpg1": "base64エンコードされた画像データ1",
            "imageJpg2": "base64エンコードされた画像データ2",
            "imageJpg3": "base64エンコードされた画像データ3",
            "imageJpg4": "base64エンコードされた画像データ4",
            "FuelData": {
                "vehicleID": "車両ID",
                "refuelAt": "給油日時",
                "location": "",
                "amountL": 19.87,
                "totalFee": 12312,
                "literFee": null,
                "odometerKilo": 1232,
                "tripmeterKilo": null,
                "fuelType": "燃料タイプ",
                "receiptImgJPG": null
            },
            "eventData": {
                "userID": "ユーザーID",
                "startAt": "イベント開始日時",
                "endAt": null,
                "isAllDay": false,
                "eventTitle": "イベントタイトル",
                "location": "",
                "eventURL": ""
            },
            "laptimeData": {
                "trackID": "トラックID",
                "layoutID": "レイアウトID",
                "vehicleID": "uuid",
                "laptimeMs": "ミリ秒単位",
                "roadCondition": "icebahn",
                "recordAt": "RFC3339"
            }
        }
```

+ Response 201 (application/json)
```json
        {
            "message": "投稿が正常に作成されました。",
        }
```

状態管理：useState, (余裕があれば下書き保存用localStorageも)

### いいねの新規付与（トグル） [POST] [/posts/{postID}/likes]🥺😵‍💫

ユーザーが特定の投稿に「いいね」をつける

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)
```json
        {
            "userID": "ユーザーID",
            "postID": "投稿のID"
        }
```

+ Response 201 (application/json)
```json
        {
            "message": "いいね操作を正常に完了しました",
        }
```

状態管理：useState

### 全ての投稿情報の取得 [GET] [/posts]🥺😵‍💫

ユーザーの投稿の情報を取得します。
タイムラインのグローバル枠に使用されます。（ユーザー増加によってフォロー枠に移行）

+ Parameters
    + limit:  (optional, string) - 取得する投稿の数。デフォルトは10。
    + last_postID: `SnowflakeID` (optional, string) - この日時以降に投稿された投稿のみを取得します。最初のリクエストには含まれません。ページネーション用。

+ Response 200 (application/json)
```json
        {
            "posts": [
                {
                    "postID": "post123", // 実際はuuid
                    "userData": { "UserBasicsのデータ" },
                    "vehicleData": { "VehicleBasicsのデータ" },
                    "text": "最初の投稿テキスト。",
                    "replyTo": null,
                    "repostTo": null,
                    "replyNum": 0,
                    "repostNum": 0,
                    "postType": "normal",
                    "imageFolderURL": null
                },
                {
                    "postID": "post124", // 実際はuuid
                    "userData": { "UserBasicsのデータ" },
                    "vehicleData": { "VehicleBasicsのデータ" },
                    "text": "次の投稿テキスト。",
                    "replyTo": null,
                    "repostTo": null,
                    "replyNum": 0,
                    "repostNum": 0,
                    "postType": "normal",
                    "postedAt": "2023-11-26T16:30:00Z",
                    "imageFolderURL": null
                }
            // この配列はlimitパラメータによって返される投稿の数が増減します
            ]
        }
```

+ Response 201 (application/json)
```json
        {
            "message": "投稿情報の取得に成功。",
        }
```
+ Response 429 (application/json)
```json
        {
            "error": "Too Many Requests",
            "message": "リクエスト数が多すぎます。しばらく待ってから再度試してください。"
        }
```

状態管理：SWR

### 特定のユーザーの投稿情報の取得 [GET] [/posts?userID={userID}]🥺

特定のユーザーの投稿の情報を取得します。
タイムラインのフォロー枠に使用されます。

+ Parameters
    + user_id: (required, string) - フォローテーブルから投稿を絞り込むため。必須。
    + limit:  (optional, string) - 取得する投稿の数。デフォルトは10。
    + last_postID: `SnowflakeID` (optional, string) - この日時以降に投稿された投稿のみを取得します。最初のリクエストには含まれません。ページネーション用。

+ Response 200 (application/json)

```json
        {
              "posts": [
                {
                    "postID": "post123", // 実際はSnowflakeID
                    "userData": { "UserBasicsのデータ" },
                    "vehicleData": { "VehicleBasicsのデータ" },
                    "text": "最初の投稿テキスト。",
                    "replyTo": null,
                    "repostTo": null,
                    "replyNum": 0,
                    "repostNum": 0,
                    "postType": "normal",
                    "imageFolderURL": null
                },
                {
                    "postID": "post123", // 実際はSnowflakeID
                    "userData": { "UserBasicsのデータ" },
                    "vehicleData": { "VehicleBasicsのデータ" },
                    "text": "最初の投稿テキスト。",
                    "replyTo": null,
                    "repostTo": null,
                    "replyNum": 0,
                    "repostNum": 0,
                    "postType": "normal",
                    "imageFolderURL": null
                }
                // この配列はlimitパラメータによって返される投稿の数が増減します
              ]
            }
```

+ Response 429 (application/json)
```json
        {
            "error": "Too Many Requests",
            "message": "リクエスト数が多すぎます。しばらく待ってから再度試してください。"
        }
```

状態管理：SWR

### フォローしている人の投稿情報の取得 [GET] [/posts?follow=true]🥺

フォローしているユーザーの投稿の情報を取得します。
タイムラインのフォロー枠に使用されます。

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Parameters
    + user_id: (required, string) - フォローテーブルから投稿を絞り込むため。必須。
    + limit:  (optional, string) - 取得する投稿の数。デフォルトは10。
    + last_postID: `SnowflakeID` (optional, string) - この日時以降に投稿された投稿のみを取得します。最初のリクエストには含まれません。ページネーション用。

+ Response 200 (application/json)

```json
        {
            "posts": [
                {
                    "postID": "post123", // 実際はSnowflakeID
                    "userData": { "UserBasicsのデータ" },
                    "vehicleData": { "VehicleBasicsのデータ" },
                    "text": "最初の投稿テキスト。",
                    "replyTo": null,
                    "repostTo": null,
                    "replyNum": 0,
                    "repostNum": 0,
                    "postType": "normal",
                    "imageFolderURL": null
                },
                {
                    "postID": "post123", // 実際はSnowflakeID
                    "userData": { "UserBasicsのデータ" },
                    "vehicleData": { "VehicleBasicsのデータ" },
                    "text": "最初の投稿テキスト。",
                    "replyTo": null,
                    "repostTo": null,
                    "replyNum": 0,
                    "repostNum": 0,
                    "postType": "normal",
                    "imageFolderURL": null
                }
            // この配列はlimitパラメータによって返される投稿の数が増減します
            ]
        }
```

+ Response 429 (application/json)
```json
        {
            "error": "Too Many Requests",
            "message": "リクエスト数が多すぎます。しばらく待ってから再度試してください。"
        }
```

状態管理：SWR

### いいねした人の一覧を取得する [GET] [/posts/{postID}/likes]🥺

特定の`post_id`に対してのいいねの一覧を取得します。
いいねのところにしようされます

+ Parameters
    + post_id: (required, string) - 投稿の特定に使用。
    + limit:  (optional, string) - 取得するいいねの数。デフォルトは30。
    + last_liked_at: `2020-01-01T00:00:00Z` (optional, string) - この日時以降につけられたいいねのみを取得します。最初のリクエストには含まれません。ページネーション用。

+ Response 200 (application/json)

```json
        {
            "likes": [
            {
                "userID": "post123", // 実際はuuid
                "likedAt": "RFC3339",
            },
            {
                "userID": "post123", // 実際はuuid
                "likedAt": "RFC3339",
            }
            // この配列はlimitパラメータによって返される投稿の数が増減します
            ],
            "newLastLikedAt": ""
        }
```

+ Response 429 (application/json)
```json
        {
            "error": "Too Many Requests",
            "message": "リクエスト数が多すぎます。しばらく待ってから再度試してください。"
        }
```
        状態管理：SWR

### 共有した人の一覧を取得する [GET] [/posts/{postID}/reposts]←共有用中間テーブルが必要

特定の`postID`に対してのリポストの一覧を取得します。

+ Parameters
    + postID: (required, string) - 投稿の特定に使用。
    + limit:  (optional, string) - 取得する投稿の数。デフォルトは10。
    + lastLikedAt: `2020-01-01T00:00:00Z` (optional, string) - この日時以降の投稿のみを取得します。最初のリクエストには含まれません。ページネーション用。

+ Response 200 (application/json)

```json
        {
            "likes": [
            {
                "userID": "post123", // 実際はuuid
                "likedAt": "RFC3339",
            },
            {
                "userID": "post123", // 実際はuuid
                "likedAt": "RFC3339",
            }
            // この配列はlimitパラメータによって返される投稿の数が増減します
            ]
        }
```

状態管理：SWR

### リプライの一覧を取得する [GET] [/posts/{postID}/replies]←リプライ用中間テーブルが必要

### 投稿情報の削除 [DELETE] [/posts/{postID}]

指定された`postID`の投稿を削除します。

+ Parameters
    + post_id: (required, string) - 投稿の特定に使用します。

+ Request (application/json)
    + Headers: `Authorization: Bearer {アクセストークン}`

+ Response 204 (application/json)

    + Body
        ```json
            {
                "message": "投稿が正常に削除されました。"
            }
        ```

+ Response 404 (application/json)

    + Body
        ```json
            {
                "error": "指定された投稿が見つかりません。"
            }
        ```
+ Response 403 (application/json)

    + Body
        ```json
            {
                "error": "この操作を行う権限がありません。"
            }
        ```
状態管理：useState

## エンドポイント [/fuels]

### 給油情報の更新 [PATCH] [/fuels/{fuelID}]

Authorizationされているユーザー、かつVehicle_User_Links内`access_level`=`"admin"`,`"moderater"`の場合、給油情報の更新を許可します。

～バックエンドの人へ～
Cassandraでどう実装するかですが、プライマリキーを完全に一致させた上で、INSERTを行えば多分実装可能。
普段の投稿に載せたfuelのPOSTの場合はIF NOT EXISTをクエリに載せますが、こちらの場合はそれを外します。

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)
```json
        {
            "refualAt": "RFC3339",
            "amountMl": "給油量(単位ml)",
            "fuelType": "油種。dieselとか",
            "literFee": "リッター当たりの値段",
            "totalFee": "今回給油にかかった金額",
            "odometerKilo": "給油時点の走行距離",
            "tripmeterKilo": "前回の給油からの距離",
            "location": "場所",
            "receiptImgJpg": "レシート画像"
        }
```
+ Response 200 (application/json)
```json
    {
        "message": "給油情報を更新しました。"
     }
```

状態管理：useState

## エンドポイント [/events]

### イベント情報の削除 [DELETE] [/events/{event_id}]

指定された`eventID`の投稿を削除します。

+ Parameters
    + eventID: (required, string) - 投稿の特定に使用します。

+ Request (application/json)
    + Headers: `Authorization: Bearer {アクセストークン}`

+ Response 204 (application/json)

    + Body
        ```json
            {
                "message": "イベントが正常に削除されました。"
            }
        ```

+ Response 404 (application/json)

    + Body
        ```json
            {
                "error": "指定されたイベントが見つかりません。"
            }
        ```
+ Response 403 (application/json)

    + Body
        ```json
            {
                "error": "この操作を行う権限がありません。"
            }
        ```

状態管理：useState


## エンドポイント [/laptimes]

### ラップタイムの更新 [PATCH] [/laptimes/{laptime_id}]

+ ヘッダー: `Authorization: Bearer {access_token}`

+ Request (application/json)
```json
    {
        "laptimeID": "指定されたラップタイムID",
        "updates": {
            "laptime": {
            "trackID": "新しいトラックID",
            "layoutID": "新しいレイアウトID",
            "vehicleID": "uuid",
            "laptimeMs": "ミリ秒単位",
            "roadCondition": "icebahn",
            "recordAt": "RFC3339"
            },
        }
    }
```
+ Response 200 (application/json)
```json
    {
        "message": "ユーザー情報を更新しました。"
    }
```

状態管理：useState

### ラップタイムの削除 [DELETE] [/laptimes/{laptime_id}]

指定された`laptimeID`の投稿を削除します。

+ Parameters
    + laptime_id: (required, string) - ラップタイムの特定に使用します。

+ Request (application/json)
    + Headers: `Authorization: Bearer {アクセストークン}`

+ Response 204 (application/json)

    + Body
        ```json
            {
                "message": "ラップタイムが正常に削除されました。"
            }
        ```

+ Response 404 (application/json)

    + Body
        ```json
            {
                "error": "指定されたラップタイムが見つかりません。"
            }
        ```
+ Response 403 (application/json)

    + Body
        ```json
            {
                "error": "この操作を行う権限がありません。"
            }
        ```

## エンドポイント [/tracks]

### トラックの新規登録 [POST] [/tracks]

+ ヘッダー: `Authorization: Bearer {access_token}`

トラックを新規作成します。

+ Request (application/json)
```json
    {
        "tracks": [
            {
                "trackName": "Nordschleife",
                "country": "Germany",
                "location": "Nurburg"
            }
        ],
        "layouts": [
            {
                "layoutName": "Endurance",
                "latitude": "",
                "longitude": "",
                "length": "25415"
            }
        ]
    }
```

+ Response 201 (application/json)
```json
        {
            "message": "トラックが正常に作成されました。",
        }
```

状態管理：useState

### レイアウトの追加 [POST] [/tracks/{track_id}/layouts]

トラックのレイアウトを新規作成します。

+ Request (application/json)
```json
        {
            "layoutName": "Endurance",
            "latitude": "",
            "longitude": "",
            "length": 25415
        }
```

+ Response 201 (application/json)
```json
        {
            "message": "レイアウトが正常に作成されました。",
        }
```
状態管理：useState

## 別のエンドポイント [/another/path]

...
