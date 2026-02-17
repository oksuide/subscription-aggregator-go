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
	DeleteSubscription(ctx context.Context, id int64) error
	UpdateSubscription(ctx context.Context, id int64, updates Subscription) (Subscription, error)
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
			return Subscription{}, ErrSubscriptionNotFound
		}
		r.logger.Error("failed to get subscription from db",
			zap.Int64("id", id),
			zap.Error(err))
		return Subscription{}, fmt.Errorf("db select: %w", err)
	}
	return received, nil
}

func (r *Repo) DeleteSubscription(ctx context.Context, id int64) error {
	r.logger.Debug("executing delete query",
		zap.Int64("subscription_id", id))

	query := "DELETE FROM subscriptions WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete subscription from db",
			zap.Int64("id", id),
			zap.Error(err))
		return fmt.Errorf("db delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("subscription not found for deletion",
			zap.Int64("id", id))
		return ErrSubscriptionNotFound
	}
	return nil
}

func (r *Repo) UpdateSubscription(ctx context.Context, id int64, updates Subscription) (Subscription, error) {
	r.logger.Debug("executing update query",
		zap.Int64("subscription_id", id))
	query := `UPDATE subscriptions 
			  SET service_name = $2, 
			  	  price = $3,
			  	  user_id = $4,
			  	  start_date = $5,
			  	  end_date = $6,
			  	  updated_at = NOW()
			  	  WHERE id = $1
				  RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`
	var updated Subscription
	err := r.db.QueryRowxContext(ctx, query,
		id,
		updates.ServiceName,
		updates.Price,
		updates.UserID,
		updates.StartDate,
		updates.EndDate,
	).StructScan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("subscription not found in db", zap.Int64("id", id))
			return Subscription{}, ErrSubscriptionNotFound
		}

		r.logger.Error("failed to update subscription in db",
			zap.Error(err),
			zap.Int64("id", id))
		return Subscription{}, fmt.Errorf("db update: %w", err)
	}
	return updated, nil
}
