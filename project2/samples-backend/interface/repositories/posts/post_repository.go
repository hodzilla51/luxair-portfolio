package repositories

import (
	"database/sql"
	models "sample-backend-go/internal/domain/entities"
	"strconv"

	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepository struct {
	pool    *pgxpool.Pool
	session *gocql.Session
}

func NewPostRepository(pool *pgxpool.Pool, session *gocql.Session) *PostRepository {
	return &PostRepository{pool: pool, session: session}
}

func (r *PostRepository) InsertPostToPostgres(postData models.Post, partitionSeed string) error {
	query := `INSERT INTO posts (
		user_id, post_id, post_type, vehicle_id, reply_to, repost_to, reply_num, repost_num, text,
		image_folder_url, partition_seed
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(context.Background(), query,
		postData.UserID, postData.PostID, postData.PostType, postData.VehicleID,
		postData.ReplyTo, postData.RepostTo, 0, 0,
		postData.Text, postData.ImageFolderURL, partitionSeed)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostRepository) InsertFuelToCassandra(fuelData models.Fuel) error {
	query := `INSERT INTO fuels (
		vehicle_id, bucket, fuel_id, post_id, refuel_at, amount_ml, fuel_type, liter_fee,
		total_fee, odometer_kilo, tripmeter_kilo, location, receipt_img_url
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if err := r.session.Query(query,
		fuelData.VehicleID, fuelData.Bucket, fuelData.FuelID, fuelData.PostID, fuelData.RefuelAt, fuelData.AmountML,
		fuelData.FuelType, fuelData.LiterFee, fuelData.TotalFee, fuelData.OdometerKilo,
		fuelData.TripmeterKilo, fuelData.Location, fuelData.ReceiptImgURL).Consistency(gocql.Quorum).Exec(); err != nil {
		return err
	}
	return nil
}
func (r *PostRepository) InsertEventToCassandra(eventData models.Event) error {
	query := `INSERT INTO events (
		user_id, bucket, event_id, start_at, end_at, event_title, is_allday,
		location, post_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if err := r.session.Query(query,
		eventData.UserID, eventData.Bucket, eventData.EventID, eventData.StartAt, eventData.EndAt,
		eventData.EventTitle, eventData.IsAllDay, eventData.Location, eventData.PostID).Consistency(gocql.Quorum).Exec(); err != nil {
		return err
	}
	return nil
}
func (r *PostRepository) InsertLaptimeDetailsToPostgres(laptimeData models.Laptime, partitionSeed string) error {
	query := `INSERT INTO laptime_details (
		laptime_id, track_id, layout_id, user_id, vehicle_id, laptime_ms, road_condition, record_at, partition_seed
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(context.Background(), query,
		laptimeData.LaptimeID, laptimeData.TrackID, laptimeData.LayoutID, laptimeData.UserID,
		laptimeData.VehicleID, laptimeData.LaptimeMs, laptimeData.RoadCondition, laptimeData.RecordAt,
		partitionSeed)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetPosts(limit int, lastPostID int64) ([]models.PostDelivery, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	var deliveries []models.PostDelivery
	query := `SELECT p.post_id, p.user_id, p.post_type, p.vehicle_id, p.reply_to, p.repost_to, p.like_num, p.reply_num, p.repost_num, p.text, p.image_folder_url, u.username, u.handlename, u.user_images_url, v.model_type, v.nickname, v.vehicle_images_url
			FROM posts p
			JOIN user_basics u ON p.user_id = u.user_id
			LEFT JOIN vehicle_basics v ON p.vehicle_id = v.vehicle_id AND p.vehicle_id IS NOT NULL
			WHERE ($1::bigint = 0 OR p.post_id < $1::bigint)
			ORDER BY p.post_id DESC
			LIMIT $2`

	args := []interface{}{lastPostID, limit} //←以降の実装でパーティションキーが必要になった場合、右側の{}に入るように実装すればいけるで

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var delivery models.PostDelivery
		var postID int64
		var userID, username, handlename, userImagesURL string
		var vehicleID, modelType, nickname, vehicleImagesURL sql.NullString // 車両情報がNULL可能性があるため、NullStringを使用
		// rows.Scanの引数リストに、車両情報に関連する変数を追加
		if err := rows.Scan(&postID, &userID, &delivery.PostType, &vehicleID, &delivery.ReplyTo, &delivery.RepostTo, &delivery.LikeNum, &delivery.ReplyNum, &delivery.RepostNum, &delivery.Text, &delivery.ImageFolderURL, &username, &handlename, &userImagesURL, &modelType, &nickname, &vehicleImagesURL); err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}

		delivery.PostID = strconv.FormatInt(postID, 10)

		delivery.UserData.UserID = userID
		delivery.UserData.Username = username
		delivery.UserData.Handlename = handlename
		delivery.UserData.UserImagesURL = userImagesURL

		delivery.VehicleData.VehicleID = vehicleID.String
		delivery.VehicleData.ModelType = modelType.String
		delivery.VehicleData.Nickname = nickname.String
		delivery.VehicleData.VehicleImagesURL = vehicleImagesURL.String
		deliveries = append(deliveries, delivery)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}
	return deliveries, nil
}
func (r *PostRepository) GetPostsUser(limit int, lastPostID int64, userID string, partitionSeed string) ([]models.PostDelivery, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	var posts []models.PostDelivery
	// クエリにJOINを加えて、UserDataとVehicleDataを取得する
	query := `
	SELECT 
		p.post_id, p.text, p.repost_to, p.reply_to, p.repost_num, p.reply_num, p.post_type, p.image_folder_url,
		ub.user_id, ub.username, ub.handlename, ub.user_images_url,
		vb.vehicle_id, vb.model_type, vb.nickname, vb.vehicle_images_url
	FROM posts p
	INNER JOIN user_basics ub ON p.user_id = ub.user_id AND p.partition_seed = ub.partition_seed
	LEFT JOIN vehicle_basics vb ON p.vehicle_id = vb.vehicle_id AND p.partition_seed = vb.partition_seed
	WHERE p.user_id = $1 AND p.partition_seed = $2`

	args := []interface{}{userID, partitionSeed}

	if lastPostID > 0 {
		query += " AND p.post_id > $3"
		args = append(args, lastPostID)
	}

	query += " ORDER BY p.post_id DESC LIMIT $4"
	args = append(args, limit)

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.PostDelivery
		var imageFolderURL sql.NullString // Optional field

		// UserDataとVehicleDataのための変数を準備
		var userID string
		var username, handlename, userImagesURL string
		var vehicleID string
		var modelType, nickname, vehicleImagesURL sql.NullString // Optional fields for VehicleData

		if err := rows.Scan(&post.PostID, &post.Text, &post.RepostTo, &post.ReplyTo, &post.RepostNum, &post.ReplyNum, &post.PostType, &imageFolderURL,
			&userID, &username, &handlename, &userImagesURL,
			&vehicleID, &modelType, &nickname, &vehicleImagesURL); err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}

		// Optional fields handling
		if imageFolderURL.Valid {
			post.ImageFolderURL = &imageFolderURL.String
		}

		// UserBasics and VehicleBasics structs filling
		post.UserData = models.UserBasics{
			UserID:        userID,
			Username:      username,
			Handlename:    handlename,
			UserImagesURL: userImagesURL,
		}

		// VehicleDataの扱いは、vehicle_idが実際に存在するかどうかで変わる
		if vehicleID != "" {
			post.VehicleData = models.VehicleBasics{
				VehicleID:        vehicleID,
				ModelType:        modelType.String, // これらのフィールドがNULL可能性があるので、適切に扱う
				Nickname:         nickname.String,
				VehicleImagesURL: vehicleImagesURL.String,
			}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return posts, nil
}
func (r *PostRepository) GetPostsFollow(limit int, lastPostID int64, userID string, partitionSeed string) ([]models.PostDelivery, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	var posts []models.PostDelivery
	// Update the query to perform the necessary joins
	query := `
    SELECT 
        p.post_id, p.text, p.reply_to, p.repost_to, p.repost_num, p.reply_num, p.post_type, p.image_folder_url,
        ub.user_id, ub.username, ub.handlename, ub.user_images_url,
        vb.vehicle_id, vb.model_type, vb.nickname, vb.vehicle_images_url
    FROM 
        posts p
        INNER JOIN followings f ON p.user_id = f.following_id AND p.partition_seed = f.partition_seed
        INNER JOIN user_basics ub ON p.user_id = ub.user_id AND p.partition_seed = ub.partition_seed
        LEFT JOIN vehicle_basics vb ON p.vehicle_id = vb.vehicle_id AND p.partition_seed = vb.partition_seed
    WHERE 
        f.user_id = $1 AND p.partition_seed = $2`

	args := []interface{}{userID, partitionSeed}

	if lastPostID > 0 {
		query += " AND p.post_id > $3"
		args = append(args, lastPostID)
	}

	query += " ORDER BY p.post_id DESC LIMIT $4"
	args = append(args, limit)

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.PostDelivery
		var imageFolderURL string // Handle optional fields

		// Prepare variables for UserBasics and VehicleBasics
		var vehicleID sql.NullString
		var modelType, nickname, vehicleImagesURL sql.NullString // Handle NULLs for VehicleData

		if err := rows.Scan(&post.PostID, &post.Text, &post.ReplyTo, &post.RepostTo, &post.RepostNum, &post.ReplyNum, &post.PostType, &imageFolderURL,
			&post.UserData.UserID, &post.UserData.Username, &post.UserData.Handlename, &post.UserData.UserImagesURL,
			&vehicleID, &modelType, &nickname, &vehicleImagesURL); err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}

		// Only fill VehicleData if vehicle_id is not empty
		if vehicleID.Valid {
			post.VehicleData.VehicleID = vehicleID.String
			post.VehicleData.ModelType = modelType.String
			post.VehicleData.Nickname = nickname.String
			post.VehicleData.VehicleImagesURL = vehicleImagesURL.String
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return posts, nil
}
