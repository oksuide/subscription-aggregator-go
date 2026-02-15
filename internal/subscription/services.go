package subscription

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidDateFormat    = errors.New("invalid date format, use MM-YYYY")
	ErrEndDateBeforeStart   = errors.New("end date cannot be before start date")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type Service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) CreateSubscription(ctx context.Context, rawSub SubscriptionReq) (Subscription, error) {
	s.logger.Info("creating subscription",
		zap.String("user_id", rawSub.UserID.String()),
		zap.String("service", rawSub.ServiceName),
		zap.Int("price", rawSub.Price),
		zap.String("start_date", rawSub.StartDate))

	startDate, err := time.Parse("01-2006", rawSub.StartDate)
	if err != nil {
		s.logger.Warn("invalid start date format",
			zap.String("user_id", rawSub.UserID.String()),
			zap.String("start_date", rawSub.StartDate),
			zap.Error(err))
		return Subscription{}, ErrInvalidDateFormat
	}

	var endDate *time.Time
	if rawSub.EndDate != nil {
		parsed, err := time.Parse("01-2006", *rawSub.EndDate)
		if err != nil {
			s.logger.Warn("invalid end date format",
				zap.String("user_id", rawSub.UserID.String()),
				zap.String("end_date", *rawSub.EndDate),
				zap.Error(err))
			return Subscription{}, ErrInvalidDateFormat
		}
		endDate = &parsed
	}

	if endDate != nil && endDate.Before(startDate) {
		s.logger.Warn("end date before start date",
			zap.String("user_id", rawSub.UserID.String()),
			zap.Time("start_date", startDate),
			zap.Time("end_date", *endDate))
		return Subscription{}, ErrEndDateBeforeStart
	}

	sub := Subscription{
		ServiceName: rawSub.ServiceName,
		Price:       rawSub.Price,
		UserID:      rawSub.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	response, err := s.repo.CreateSubscription(ctx, sub)
	if err != nil {
		return Subscription{}, err
	}

	return response, nil
}
