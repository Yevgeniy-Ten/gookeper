package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"gokeeper/internal/models"
)

func (p *Postgres) CreateSecret(ctx context.Context, userID int64, dto *models.CreateSecretDTO) (*models.Secret, error) {
	meta, _ := json.Marshal(dto.Meta)
	query := `
		INSERT INTO secrets (user_id, type, data, meta)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, type, data, meta, created_at, updated_at`
	var s models.Secret
	if err := p.db.QueryRowContext(ctx, query, userID, dto.Type, dto.Data, meta).Scan(
		&s.ID, &s.UserID, &s.Type, &s.Data, &s.Meta, &s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (p *Postgres) GetSecret(ctx context.Context, userID, secretID int64) (*models.Secret, error) {
	query := `
		SELECT id, user_id, type, data, meta, created_at, updated_at
		FROM secrets WHERE id = $1 AND user_id = $2`
	var s models.Secret
	if err := p.db.QueryRowContext(ctx, query, secretID, userID).Scan(
		&s.ID, &s.UserID, &s.Type, &s.Data, &s.Meta, &s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (p *Postgres) ListSecrets(ctx context.Context, userID int64) ([]*models.Secret, error) {
	query := `
		SELECT id, user_id, type, data, meta, created_at, updated_at
		FROM secrets WHERE user_id = $1 ORDER BY id DESC`
	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var secrets []*models.Secret
	for rows.Next() {
		var s models.Secret
		if err := rows.Scan(&s.ID, &s.UserID, &s.Type, &s.Data, &s.Meta, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		secrets = append(secrets, &s)
	}
	return secrets, nil
}

func (p *Postgres) DeleteSecret(ctx context.Context, userID, secretID int64) error {
	query := `DELETE FROM secrets WHERE id = $1 AND user_id = $2`
	_, err := p.db.ExecContext(ctx, query, secretID, userID)
	return err
}
