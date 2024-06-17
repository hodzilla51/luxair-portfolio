package repositories

import (
	"context"
	"fmt"
	"time"

	models "sample-backend-go/internal/domain/entities"
	responses "sample-backend-go/internal/usecase/models/responses"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VehicleRepository struct {
	pool    *pgxpool.Pool
	session *gocql.Session
}

func NewVehicleRepository(pool *pgxpool.Pool, session *gocql.Session) *VehicleRepository {
	return &VehicleRepository{pool: pool, session: session}
}

func (r *VehicleRepository) AddVehicleData(vehicleBasics models.VehicleBasicsPg, vehicleDetails models.VehicleDetailsPg, vehicleUserLink models.VehicleUserLinkPg, userID string) error {
	ctx := context.Background()
	// Start PostgreSQL transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting PostgreSQL transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // Re-throw panic
		} else if err != nil {
			tx.Rollback(ctx) // Rollback on error
		}
	}()

	// Insert data into vehicle_basics
	if err := r.AddVehicleBasics(tx, vehicleBasics); err != nil {
		return err
	}

	// Insert data into vehicle_details
	if err := r.AddVehicleDetails(tx, vehicleDetails); err != nil {
		return err
	}

	// Insert data into vehicle_user_links
	if err := r.AddVehicleUserLink(tx, vehicleUserLink, userID); err != nil {
		return err
	}

	// Commit PostgreSQL transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing PostgreSQL transaction: %w", err)
	}

	return nil
}
func (r *VehicleRepository) AddVehicleBasics(tx pgx.Tx, data models.VehicleBasicsPg) error {
	query := `INSERT INTO vehicle_basics (vehicle_id, model_type, nickname, vehicle_images_url, partition_seed) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(context.Background(), query, data.VehicleID, data.ModelType, data.Nickname, data.VehicleImagesURL, data.PartitionSeed)
	if err != nil {
		return fmt.Errorf("error inserting into vehicle_basics: %w", err)
	}
	return nil
}
func (r *VehicleRepository) AddVehicleDetails(tx pgx.Tx, data models.VehicleDetailsPg) error {
	query := `INSERT INTO vehicle_details (vehicle_id, joined_at, manufacture_date, newcar_registration_date, status, mileage, fan_num, vehicle_bio_text, frame_no, partition_seed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := tx.Exec(context.Background(), query, data.VehicleID, data.JoinedAt, data.ManufactureDate, data.NewCarRegistrationDate, data.Status, data.Mileage, 0, data.VehicleBioText, data.FrameNo, data.PartitionSeed)
	if err != nil {
		return fmt.Errorf("error inserting into vehicle_details: %w", err)
	}
	return nil
}
func (r *VehicleRepository) AddVehicleUserLink(tx pgx.Tx, link models.VehicleUserLinkPg, userID string) error {
	query := `INSERT INTO vehicle_user_links (user_id, vehicle_id, relation_type, start_at, price, created_at, access_level, transfer_type, mileage_at_registration, partition_seed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := tx.Exec(context.Background(), query, userID, link.VehicleID, link.RelationType, link.StartAt, link.Price, link.CreatedAt, link.AccessLevel, link.TransferType, link.MileageAt, link.PartitionSeed)
	if err != nil {
		return fmt.Errorf("error inserting into vehicle_user_links: %w", err)
	}
	return nil
}

// GET
func (r *VehicleRepository) GetVehicleBasic(vehicleID string, partitionSeed string) (models.VehicleBasics, error) {
	var data models.VehicleBasics

	query := `SELECT vehicle_id, model_type, nickname, vehicle_images_url 
			FROM vehicle_basics 
			WHERE vehicle_id = $1 AND partition_seed = $2`
	err := r.pool.QueryRow(context.Background(), query, vehicleID, partitionSeed).Scan(
		&data.VehicleID,
		&data.ModelType,
		&data.Nickname,
		&data.VehicleImagesURL,
	)
	if err != nil {
		return models.VehicleBasics{}, fmt.Errorf("error fetching vehicle_basics: %w", err)
	}
	return data, nil
}
func (r *VehicleRepository) GetVehicleDetail(vehicleID string, partitionSeed string) (models.VehicleDetails, error) {
	var data models.VehicleDetails

	query := `SELECT joined_at, manufacture_date, newcar_registration_date, status, mileage, fan_num, vehicle_bio_text, frame_no, partition_seed 
			FROM vehicle_details 
			WHERE vehicle_id = $1 AND partition_seed = $2`
	err := r.pool.QueryRow(context.Background(), query, vehicleID, partitionSeed).Scan(
		&data.JoinedAt,
		&data.ManufactureDate,
		&data.NewCarRegistrationDate,
		&data.Status,
		&data.Mileage,
		&data.FanNum,
		&data.VehicleBioText,
		&data.FrameNo,
		&data.PartitionSeed,
	)
	if err != nil {
		return models.VehicleDetails{}, fmt.Errorf("error fetching vehicle_details: %w", err)
	}
	return data, nil
}
func (r *VehicleRepository) GetVehicleUsers(vehicleID string) ([]responses.UserInfosResponse, error) {
	var userVehicles []responses.UserInfosResponse

	query := `
		SELECT
		ub.user_id,
		ub.username,
		ub.handlename,
		ub.user_images_url,
		vul.relation_type,
		vul.start_at,
		vul.end_at,
		vul.created_at,
		vul.access_level,
		vul.transfer_type,
		vul.mileage_at_registration
		FROM
			user_basics ub
		JOIN
			vehicle_user_links vul ON ub.user_id = vul.user_id
		WHERE
			vul.vehicle_id = $1
    `

	rows, err := r.pool.Query(context.Background(), query, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("querying vehicle info for user %v: %w", vehicleID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var v responses.UserInfosResponse
		err := rows.Scan(&v.UserID, &v.Username, &v.Handlename, &v.UserImagesURL, &v.RelationType, &v.StartAt, &v.EndAt, &v.CreatedAt, &v.AccessLevel, &v.TransferType, &v.MileageAtRegistration)
		if err != nil {
			return nil, fmt.Errorf("scanning vehicle info for user %v: %w", vehicleID, err)
		}
		userVehicles = append(userVehicles, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over vehicle info rows for user %v: %w", vehicleID, err)
	}

	return userVehicles, nil
}
func (r *VehicleRepository) GetVehicleFuels(vehicleID string, lastRefuelAt time.Time, limit int, bucket int) ([]models.Fuel, error) {
	fuels := []models.Fuel{}

	var query string
	var iter *gocql.Iter

	// lastRefuelAtがゼロ値でない場合（前回の給油日時が指定されている場合）、それを基準にデータを取得
	if !lastRefuelAt.IsZero() {
		query = `SELECT vehicle_id, bucket, fuel_id, post_id, refuel_at, amount_ml, fuel_type, liter_fee, total_fee, odometer_kilo, tripmeter_kilo, location, receipt_img_url FROM fuels WHERE vehicle_id = ? AND bucket = ? AND refuel_at < ? ORDER BY refuel_at DESC LIMIT ?`
		iter = r.session.Query(query, vehicleID, bucket, lastRefuelAt, limit).Iter()
	} else {
		// lastRefuelAtがゼロ値の場合（最初のページ）、最新のデータから取得
		query = `SELECT vehicle_id, bucket, fuel_id, post_id, refuel_at, amount_ml, fuel_type, liter_fee, total_fee, odometer_kilo, tripmeter_kilo, location, receipt_img_url FROM fuels WHERE vehicle_id = ? AND bucket = ? ORDER BY refuel_at DESC LIMIT ?`
		iter = r.session.Query(query, vehicleID, bucket, limit).Iter()
	}

	var fuel models.Fuel
	for iter.Scan(&fuel.VehicleID, &fuel.Bucket, &fuel.FuelID, &fuel.PostID, &fuel.RefuelAt, &fuel.AmountML, &fuel.FuelType, &fuel.LiterFee, &fuel.TotalFee, &fuel.OdometerKilo, &fuel.TripmeterKilo, &fuel.Location, &fuel.ReceiptImgURL) {
		fuels = append(fuels, fuel)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error iterating fuels records: %w", err)
	}

	return fuels, nil
}
func (r *VehicleRepository) GetVehicleLaptimes(vehicleID string, partitionSeed string) ([]models.Laptime, error) {
	var laptimes []models.Laptime

	query := `
    SELECT *
    FROM laptime_details
    WHERE user_id = $1 AND partition_seed = $2
    `

	rows, err := r.pool.Query(context.Background(), query, vehicleID, partitionSeed)
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
