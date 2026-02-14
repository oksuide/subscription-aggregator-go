package subscription

import "go.uber.org/zap"

type Service struct {
	Repo   Repository
	Logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{
		Repo:   repo,
		Logger: logger,
	}
}
