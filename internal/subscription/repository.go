package subscription

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

type Repository interface {
	CreateSubscription(ctx context.Context, sub Subscription) (Subscription, error)
	GetSubscription(ctx context.Context, id int64) (Subscription, error)
}

func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &Repo{
		db:     db,
		logger: logger,
	}
}

func (r *Repo) CreateSubscription(ctx context.Context, sub Subscription) (Subscription, error) {
	query := `
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
	var created Subscription
	created.ServiceName = sub.ServiceName
	created.Price = sub.Price
	created.UserID = sub.UserID
	created.StartDate = sub.StartDate
	created.EndDate = sub.EndDate
	err := r.db.QueryRowContext(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(&created.ID, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		r.logger.Error("failed to create subscription in db",
			zap.Error(err),
			zap.String("user_id", sub.UserID.String()))
		return Subscription{}, fmt.Errorf("db insert: %w", err)
	}
	return created, nil
}

func (r *Repo) GetSubscription(ctx context.Context, id int64) (Subscription, error) {
	query := "SELECT * FROM subscriptions WHERE id = $1"

	var received Subscription

	err := r.db.GetContext(ctx, &received, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("subscription not found in db", zap.Int64("id", id))
			return Subscription{}, sql.ErrNoRows
		}
		r.logger.Error("failed to get subscription from db",
			zap.Int64("id", id),
			zap.Error(err))
		return Subscription{}, fmt.Errorf("db select: %w", err)
	}
	return received, nil
}
