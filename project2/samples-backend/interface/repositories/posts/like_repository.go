package repositories

import (
	"context"
	"log"
	"sample-backend-go/internal/usecase/models/responses"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type LikeRepository struct {
	pool *pgxpool.Pool
}

func NewLikeRepository(pool *pgxpool.Pool) *LikeRepository {
	return &LikeRepository{pool: pool}
}

// LikeExists - 特定のいいねが存在するかチェック
func (r *LikeRepository) LikeExists(postID int64, userID string, partitionSeed string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2 AND partition_seed = $3)`, postID, userID, partitionSeed).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// DeleteLike - いいねを削除
func (r *LikeRepository) DeleteLike(postID int64, userID string, partitionSeedLike string, partitionSeedPost string) error {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM likes WHERE post_id = $1 AND user_id = $2 AND partition_seed = $3`, postID, userID, partitionSeedLike)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	_, err = tx.Exec(context.Background(), `UPDATE posts SET like_num = like_num - 1 WHERE post_id = $1 AND partition_seed = $2`, postID, partitionSeedPost)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// AddLike - いいねを追加
func (r *LikeRepository) AddLike(postID int64, userID string, partitionSeedLike string, partitionSeedPost string) error {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	now := time.Now()

	_, err = tx.Exec(context.Background(), `INSERT INTO likes (post_id, user_id, liked_at, partition_seed) VALUES ($1, $2, $3, $4)`, postID, userID, now, partitionSeedLike)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	log.Println(postID)
	log.Println(partitionSeedPost)
	_, err = tx.Exec(context.Background(), `UPDATE posts SET like_num = like_num + 1 WHERE post_id = $1 AND partition_seed = $2`, postID, partitionSeedPost)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// GetLikesWithPagination - いいねしたユーザーのIDと最も古いいいねの時間をページネーション付きで取得
func (r *LikeRepository) GetLikesWithPagination(postID int64, lastLikedAt time.Time, limit int, partitionSeed string) ([]responses.LikeDetail, time.Time, error) {
	var likes []responses.LikeDetail
	var oldestLikedAt time.Time

	query := `SELECT user_id, liked_at FROM likes WHERE post_id = $1 AND liked_at < $2 AND partition_seed = $3 ORDER BY liked_at DESC LIMIT $4`
	rows, err := r.pool.Query(context.Background(), query, postID, lastLikedAt, partitionSeed, limit)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var detail responses.LikeDetail
		if err := rows.Scan(&detail.UserID, &detail.LikedAt); err != nil {
			return nil, time.Time{}, err
		}
		likes = append(likes, detail)
		if oldestLikedAt.IsZero() || detail.LikedAt.Before(oldestLikedAt) {
			oldestLikedAt = detail.LikedAt
		}
	}

	if err := rows.Err(); err != nil {
		return nil, time.Time{}, err
	}

	// 最後に取得したliked_atが最も古いいいねの時間になる場合のロジックは調整が必要です。
	// 最後に取得したレコードのlikedAtをそのまま使うか、別の方法で特定する必要があります。
	return likes, oldestLikedAt, nil
}
