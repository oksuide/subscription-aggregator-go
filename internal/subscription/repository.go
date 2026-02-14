package subscription

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repo struct {
	DB     *sqlx.DB
	Logger *zap.Logger
}

type Repository interface {
	// IsUserExistWithEmail(email string) (bool, error)
}

func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &Repo{
		DB:     db,
		Logger: logger,
	}
}
