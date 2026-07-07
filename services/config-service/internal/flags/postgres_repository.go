package flags

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(flag Flag) error {
	targetingRulesJSON, err := json.Marshal(flag.TargetingRules)
	if err != nil {
		return err
	}

	id, err := newUUID()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO flags (
			id,
			key,
			name,
			description, 
			enabled,
			rollout_percentage,
			targeting_rules,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9)
	`,
		id,
		flag.Key,
		flag.Name,
		flag.Description,
		flag.Enabled,
		flag.RolloutPercentage,
		string(targetingRulesJSON),
		flag.CreatedAt,
		flag.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrFlagAlreadyExists
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) List() ([]Flag, error) {
	rows, err := r.db.Query(`
		SELECT
			key,
			name,
			description,
			enabled,
			rollout_percentage,
			targeting_rules,
			created_at,
			updated_at
		FROM flags
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flags := []Flag{}

	for rows.Next() {
		var flag Flag
		var targetingRulesJSON []byte

		if err := rows.Scan(
			&flag.Key,
			&flag.Name,
			&flag.Description,
			&flag.Enabled,
			&flag.RolloutPercentage,
			&targetingRulesJSON,
			&flag.CreatedAt,
			&flag.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(targetingRulesJSON, &flag.TargetingRules); err != nil {
			return nil, err
		}

		flags = append(flags, flag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return flags, nil
}

func (r *PostgresRepository) GetByKey(key string) (Flag, error) {
	var flag Flag
	var targetingRulesJSON []byte

	err := r.db.QueryRow(`
		SELECT
			key,
			name,
			description,
			enabled,
			rollout_percentage,
			targeting_rules,
			created_at,
			updated_at
		FROM flags
		WHERE key = $1
	`, key).Scan(
		&flag.Key,
		&flag.Name,
		&flag.Description,
		&flag.Enabled,
		&flag.RolloutPercentage,
		&targetingRulesJSON,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Flag{}, ErrFlagNotFound
		}

		return Flag{}, err
	}

	if err := json.Unmarshal(targetingRulesJSON, &flag.TargetingRules); err != nil {
		return Flag{}, err
	}

	return flag, nil
}

func (r *PostgresRepository) Update(flag Flag) error {
	targetingRulesJSON, err := json.Marshal(flag.TargetingRules)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(`
		UPDATE flags
		SET
			name = $1,
			description = $2,
			enabled = $3,
			rollout_percentage = $4,
			targeting_rules = $5::jsonb,
			updated_at = $6
		WHERE key = $7
		`,
		flag.Name,
		flag.Description,
		flag.Enabled,
		flag.RolloutPercentage,
		string(targetingRulesJSON),
		flag.UpdatedAt,
		flag.Key,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrFlagNotFound
	}

	return nil
}

func newUUID() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:],
	), nil
}

func (r *PostgresRepository) RecordExposure(event ExposureEvent) error {
	id, err := newUUID()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO exposure_events (
			id,
			flag_key,
			user_id,
			enabled,
			reason,
			bucket,
			rollout_percentage,
			created_at
		)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		id,
		event.FlagKey,
		event.UserID,
		event.Enabled,
		event.Reason,
		event.Bucket,
		event.RolloutPercentage,
		event.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetExposureSummary(flagKey string) (ExposureSummary, error) {
	summary := ExposureSummary{
		FlagKey: flagKey,
	}

	err := r.db.QueryRow(`
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE enabled = true),
			COUNT(*) FILTER (WHERE enabled = false)
		FROM exposure_events
		WHERE flag_key = $1
	`, flagKey).Scan(
		&summary.Total,
		&summary.Enabled,
		&summary.Disabled,
	)
	if err != nil {
		return ExposureSummary{}, err
	}
	return summary, nil
}

func (r *PostgresRepository) RecordConversion(event ConversionEvent) error {
	id, err := newUUID()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO conversion_events (
			id,
			flag_key,
			user_id,
			event_name,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`,
		id,
		event.FlagKey,
		event.UserID,
		event.EventName,
		event.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}