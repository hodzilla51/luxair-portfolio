package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	models "sample-backend-go/internal/domain/entities"
	"sample-backend-go/internal/domain/entities/config"
	"sample-backend-go/internal/domain/services/keycloak"
	commonUtils "sample-backend-go/internal/domain/services/utils"
	"sample-backend-go/internal/usecase/models/requests"
	"sample-backend-go/internal/usecase/models/responses"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// POST
func (r *UserRepository) AddUserData(basic models.UserBasicsPg, detail models.UserDetailsPg) error {

	// トランザクションを開始
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(context.Background())
			panic(p) // パニックを再スロー
		} else if err != nil {
			tx.Rollback(context.Background()) // エラーがあればロールバック
		}
	}()
	err = r.addUserBasic(tx, basic)
	if err != nil {
		return fmt.Errorf("AddUserBasic() failed: %w", err)
	}

	err = r.addUserDetail(tx, detail)
	if err != nil {
		return fmt.Errorf("AddUserDetail() failed: %w", err)
	}

	// トランザクションをコミット
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
func (r *UserRepository) addUserBasic(tx pgx.Tx, data models.UserBasicsPg) error {
	query := `INSERT INTO user_basics (user_id, username, handlename, user_images_url, partition_seed) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(context.Background(), query, data.UserID, data.Username, data.Handlename, "", data.PartitionSeed)
	if err != nil {
		return fmt.Errorf("error inserting into user_basics: %w", err)
	}
	return nil
}
func (r *UserRepository) addUserDetail(tx pgx.Tx, data models.UserDetailsPg) error {
	query := `INSERT INTO user_details (user_id, joined_at, following_num, follower_num, location, simulator, user_bio_text, partition_seed
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := tx.Exec(context.Background(), query, data.UserID, data.JoinedAt, 0, 0, "", "", "", data.PartitionSeed)
	if err != nil {
		return fmt.Errorf("error inserting into users_details: %w", err)
	}
	return nil
}

// GET
func (r *UserRepository) GetUserBasic(userID string) (models.UserBasics, error) {
	var data models.UserBasics

	partitionSeed, err := commonUtils.GenerateSeedUUID256(userID)
	if err != nil {
		log.Printf("パーティショニング用シード値の生成に失敗: %v", err)
		return models.UserBasics{}, err
	}
	query := `SELECT user_id, username, handlename, user_images_url FROM user_basics WHERE user_id = $1 AND partition_seed = $2`
	err = r.pool.QueryRow(context.Background(), query, userID, partitionSeed).Scan(&data.UserID, &data.Username, &data.Handlename, &data.UserImagesURL)
	if err != nil {
		return models.UserBasics{}, fmt.Errorf("error fetching user_basics: %w", err)
	}
	return data, nil
}
func (r *UserRepository) GetUserDetail(userID string) (models.UserDetails, error) {
	var data models.UserDetails

	partitionSeed, err := commonUtils.GenerateSeedUUID256(userID)
	if err != nil {
		log.Printf("パーティショニング用シード値の生成に失敗: %v", err)
		return models.UserDetails{}, err
	}
	query := `SELECT user_id, joined_at, following_num, follower_num, location, simulator, user_bio_text FROM user_details WHERE user_id = $1 AND partition_seed = $2`
	err = r.pool.QueryRow(context.Background(), query, userID, partitionSeed).Scan(&data.UserID, &data.JoinedAt, &data.FollowingNum, &data.FollowerNum, &data.Location, &data.Simulator, &data.UserBioText)
	if err != nil {
		return models.UserDetails{}, fmt.Errorf("error fetching users_details: %w", err)
	}
	return data, nil
}
func (r *UserRepository) GetUserBasicWithUsername(username string) (models.UserBasics, error) {
	var data models.UserBasics

	query := `SELECT user_id, username, handlename, user_images_url FROM user_basics WHERE username = $1`
	err := r.pool.QueryRow(context.Background(), query, username).Scan(&data.UserID, &data.Username, &data.Handlename, &data.UserImagesURL)
	if err != nil {
		log.Printf("error fetching user_basics: %v", err)
		return models.UserBasics{}, fmt.Errorf("error fetching user_basics: %w", err)
	}
	return data, nil
}
func (r *UserRepository) GetUserDetailWithUsername(username string) (models.UserDetails, error) {
	var data models.UserDetails

	query := `SELECT user_id, joined_at, following_num, follower_num, location, simulator, user_bio_text FROM user_details WHERE username = $1`
	err := r.pool.QueryRow(context.Background(), query, username).Scan(&data.UserID, &data.JoinedAt, &data.FollowingNum, &data.FollowerNum, &data.Location, &data.Simulator, &data.UserBioText)
	if err != nil {
		log.Printf("error fetching user_details: %v", err)
		return models.UserDetails{}, fmt.Errorf("error fetching users_details: %w", err)
	}
	return data, nil
}
func (r *UserRepository) GetUserVehicles(userID string) ([]responses.VehicleCards, error) {
	vehicles := []responses.VehicleCards{}

	query := `
	SELECT
		v.vehicle_id AS "vehicleID",
		v.model_type,
		v.nickname,
		v.vehicle_images_url,
		vur.relation_type
	FROM
		vehicle_basics v
	JOIN
		vehicle_user_links vur ON v.vehicle_id = vur.vehicle_id
	WHERE
		vur.user_id = $1
	`

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("querying vehicles for user %v: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var v responses.VehicleCards
		if err := rows.Scan(&v.VehicleID, &v.ModelType, &v.Nickname, &v.VehicleImagesURL, &v.RelationType); err != nil {
			return nil, fmt.Errorf("scanning vehicle for user %v: %w", userID, err)
		}
		vehicles = append(vehicles, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over vehicle rows for user %v: %w", userID, err)
	}

	return vehicles, nil
}
func (r *UserRepository) GetUserLaptimes(userID string) ([]models.Laptime, error) {
	var laptimes []models.Laptime

	query := `
    SELECT *
    FROM laptime_details
    WHERE user_id = $1
    `

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var laptime models.Laptime
		err := rows.Scan(&laptime.LaptimeID, &laptime.TrackID, &laptime.LayoutID,
			&laptime.UserID, &laptime.VehicleID, &laptime.LaptimeMs,
			&laptime.RoadCondition, &laptime.RecordAt)
		if err != nil {
			return nil, err
		}
		laptimes = append(laptimes, laptime)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return laptimes, nil
}

// PATCH
func (r *UserRepository) PatchUserData(userID string, updateReq requests.UserUpdateRequest, partitionSeed string) error {
	// user_basicsテーブルの更新
	basicsQuery := `
	UPDATE user_basics
	SET handlename = $1
	WHERE user_id = $2 AND partition_seed = $3
	`

	// user_detailsテーブルの更新
	detailsQuery := `
	UPDATE user_details
	SET location = $1, user_bio_text = $2
	WHERE user_id = $3 AND partition_seed = $4
	`

	// トランザクションを開始
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// user_basicsの更新
	if _, err := tx.Exec(context.Background(), basicsQuery, updateReq.HandleName, userID, partitionSeed); err != nil {
		return fmt.Errorf("updating user_basics: %w", err)
	}

	// user_detailsの更新
	if _, err := tx.Exec(context.Background(), detailsQuery, updateReq.Location, updateReq.UserBioText, userID, partitionSeed); err != nil {
		return fmt.Errorf("updating user_details: %w", err)
	}

	// トランザクションをコミット
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
func (r *UserRepository) PatchUserProfile(userID string, userImagesURL string, req requests.UserUpdateProfileRequest, partitionSeed string) error {
	// user_basicsテーブルの更新
	basicsQuery := `
	UPDATE user_basics
	SET user_images_url = $1
	WHERE user_id = $2 AND partition_seed = $3
	`

	// user_detailsテーブルの更新
	detailsQuery := `
	UPDATE user_details
	SET user_bio_text = $1
	WHERE user_id = $2 AND partition_seed = $3
	`

	// トランザクションを開始
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// user_basicsの更新
	if _, err := tx.Exec(context.Background(), basicsQuery, userImagesURL, userID, partitionSeed); err != nil {
		return fmt.Errorf("updating user_basics: %w", err)
	}

	// user_detailsの更新
	if _, err := tx.Exec(context.Background(), detailsQuery, req.UserBioText, userID, partitionSeed); err != nil {
		return fmt.Errorf("updating user_details: %w", err)
	}

	// トランザクションをコミット
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// Keycloakの管理APIへGETリクエストを行い、入力されたusernameが存在するかをチェック
func (r *UserRepository) CheckExistsUsername(username string) (bool, error) {

	//環境変数を読み込み
	cfg, err := config.NewConfig()
	if err != nil {
		return false, err
	}

	// 管理者トークンの取得などの初期化処理をここで行う
	adminToken, err := keycloak.GetAdminToken()
	if err != nil {
		return false, err
	}
	// Keycloakのユーザー検索エンドポイントURLを構築
	url := fmt.Sprintf("%s/admin/realms/%s/users?username=%s&exact=true&lower=true", cfg.KeycloakApiUrl, cfg.KeycloakRealm, username)

	// HTTPリクエストを作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("creating request failed: %w", err)
	}

	// AuthorizationヘッダにBearerトークンを追加
	req.Header.Add("Authorization", "Bearer "+adminToken)

	// HTTPクライアントを作成してリクエストを実行
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("request to Keycloak failed: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading response body failed: %w", err)
	}
	// ユーザーが存在するかどうかをレスポンスコードで判断
	// Keycloakはユーザーが見つからない場合には200 OKと空の配列を返す
	if resp.StatusCode == http.StatusOK && string(body) != "[]" {
		return true, nil
	}
	return false, err
}

// Keycloakの管理APIへGETリクエストを行い、入力されたemailが存在するかをチェック
func (r *UserRepository) CheckExistsEmail(email string) (bool, error) {

	//環境変数を読み込み
	cfg, err := config.NewConfig()
	if err != nil {
		return false, err
	}

	// 管理者トークンの取得などの初期化処理をここで行う
	adminToken, err := keycloak.GetAdminToken()
	if err != nil {
		return false, err
	}
	// Keycloakのメールアドレス検索エンドポイントURLを構築
	url := fmt.Sprintf("%s/admin/realms/%s/users?email=%s&exact=true", cfg.KeycloakApiUrl, cfg.KeycloakRealm, url.QueryEscape(email))

	// HTTPリクエストを作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("creating request failed: %w", err)
	}

	// AuthorizationヘッダにBearerトークンを追加
	req.Header.Add("Authorization", "Bearer "+adminToken)

	// HTTPクライアントを作成してリクエストを実行
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("request to Keycloak failed: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading response body failed: %w", err)
	}

	// ユーザーが存在するかどうかをレスポンスコードで判断
	// Keycloakはユーザーが見つからない場合には200 OKと空の配列を返す
	if resp.StatusCode == http.StatusOK && string(body) != "[]" {
		return true, nil
	}
	return false, nil
}

// KeyCloakの管理APIへのPOSTリクエストを行い、新しくユーザー登録をする関数
func RegisterUserToKeycloak(userData models.UserDataForKeycloak) (string, error) {

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("設定の読み込みに失敗しました:", err)
		return "", err
	}

	adminToken, err := keycloak.GetAdminToken()
	if err != nil {
		fmt.Println("管理者トークンの取得に失敗(´；ω；｀)ﾌﾞﾜｯ", err)
		return "", err
	}
	url := fmt.Sprintf("%s/admin/realms/%s/users", cfg.KeycloakApiUrl, cfg.KeycloakRealm)

	userJSON, err := json.Marshal(userData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer"+" "+adminToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		locationHeader := resp.Header.Get("Location")
		// LocationヘッダーからIDを抽出
		userId := keycloak.ExtractUserIdFromLocation(locationHeader)
		return userId, nil
	} else {
		// レスポンスボディからエラーメッセージを読み取る
		responseBody, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("failed to register user: %s, status code: %d", responseBody, resp.StatusCode)
		return "", fmt.Errorf(errorMsg)
	}
}
