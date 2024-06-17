package repositories

import (
	"context"
	"fmt"
	models "sample-backend-go/internal/domain/entities"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type FollowRepository struct {
	pool *pgxpool.Pool
}

func NewFollowRepository(pool *pgxpool.Pool) *FollowRepository {
	return &FollowRepository{pool: pool}
}

func (r *FollowRepository) FollowExists(followingID string, userID string, partitionSeed string) (bool, error) {
	var existingFollowingID string
	err := r.pool.QueryRow(context.Background(), `SELECT following_id FROM followings WHERE following_id = $1 AND user_id = $2 AND partition_seed = $3 LIMIT 1`, followingID, userID, partitionSeed).Scan(&existingFollowingID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *FollowRepository) AddFollow(followerData models.UserFollower, followingData models.UserFollowing) error {

	// followings テーブルにデータを追加
	_, err := r.pool.Exec(context.Background(), `INSERT INTO followings (user_id, following_id, followed_at, partition_seed) VALUES ($1, $2, $3, $4)`,
		followingData.UserID, followingData.FollowingID, followingData.FollowedAt, followingData.PartitionSeed)
	if err != nil {
		return err
	}

	// followers テーブルにデータを追加
	_, err = r.pool.Exec(context.Background(), `INSERT INTO followers (user_id, follower_id, followed_at, partition_seed) VALUES ($1, $2, $3, $4)`,
		followerData.UserID, followerData.FollowerID, followerData.FollowedAt, followerData.PartitionSeed)
	if err != nil {
		return err
	}

	return nil
}

func (r *FollowRepository) DeleteFollow(followerID string, followingID string, partitionSeed string) error {

	// followings テーブルからフォロー関係を削除
	_, err := r.pool.Exec(context.Background(), `DELETE FROM followings WHERE user_id = $1 AND following_id = $2 AND partition_seed = $3`,
		followerID, followingID, partitionSeed)
	if err != nil {
		return err
	}

	// followers テーブルからフォロー関係を削除
	_, err = r.pool.Exec(context.Background(), `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2 AND partition_seed = $3`,
		followingID, followerID, partitionSeed)
	if err != nil {
		return err
	}

	return nil
}

func (r *FollowRepository) GetUserFollowings(userID string, lastFollowingAt time.Time, limit int, partitionSeed string) ([]models.UserFollowing, error) {
	if limit <= 0 {
		limit = 30 // デフォルトのリミット値
	}

	// lastFollowingAtがゼロ値でない場合のクエリを用意
	var query string
	if !lastFollowingAt.IsZero() {
		query = `
		SELECT user_id, following_id, followed_at
		FROM followings
		WHERE user_id = $1 AND partition_seed = $3 AND followed_at > $4
		ORDER BY followed_at DESC
		LIMIT $2
		`
	} else {
		query = `
		SELECT user_id, following_id, followed_at
		FROM followings
		WHERE user_id = $1 AND partition_seed = $3
		ORDER BY followed_at DESC
		LIMIT $2
		`
	}

	var rows pgx.Rows
	var err error

	if !lastFollowingAt.IsZero() {
		rows, err = r.pool.Query(context.Background(), query, userID, limit, partitionSeed, lastFollowingAt)
	} else {
		rows, err = r.pool.Query(context.Background(), query, userID, limit, partitionSeed)
	}

	if err != nil {
		return nil, fmt.Errorf("querying followings: %w", err)
	}
	defer rows.Close()

	var followings []models.UserFollowing

	for rows.Next() {
		var following models.UserFollowing
		if err := rows.Scan(&following.UserID, &following.FollowingID, &following.FollowedAt); err != nil {
			return nil, fmt.Errorf("scanning following: %w", err)
		}
		followings = append(followings, following)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return followings, nil
}
func (r *FollowRepository) GetUserFollowers(userID string, lastFollowingAt time.Time, limit int, partitionSeed string) ([]models.UserFollower, error) {
	if limit <= 0 {
		limit = 30 // デフォルトのリミット値
	}

	// lastFollowingAtがゼロ値でない場合のクエリを用意
	var query string
	if !lastFollowingAt.IsZero() {
		query = `
		SELECT user_id, follower_id, followed_at
		FROM followers
		WHERE user_id = $1 AND partition_seed = $3 AND followed_at > $4
		ORDER BY followed_at DESC
		LIMIT $2
		`
	} else {
		query = `
		SELECT user_id, follower_id, followed_at
		FROM followers
		WHERE user_id = $1 AND partition_seed = $3
		ORDER BY followed_at DESC
		LIMIT $2
		`
	}

	var rows pgx.Rows
	var err error

	if !lastFollowingAt.IsZero() {
		rows, err = r.pool.Query(context.Background(), query, userID, limit, partitionSeed, lastFollowingAt)
	} else {
		rows, err = r.pool.Query(context.Background(), query, userID, limit, partitionSeed)
	}

	if err != nil {
		return nil, fmt.Errorf("querying followers: %w", err)
	}
	defer rows.Close()

	var followings []models.UserFollower

	for rows.Next() {
		var follower models.UserFollower
		if err := rows.Scan(&follower.UserID, &follower.FollowerID, &follower.FollowedAt); err != nil {
			return nil, fmt.Errorf("scanning follower: %w", err)
		}
		followings = append(followings, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return followings, nil
}
