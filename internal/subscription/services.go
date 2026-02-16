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

func (s *Service) CreateSubscription(ctx context.Context, req SubscriptionReq) (SubscriptionResp, error) {
	s.logger.Info("creating subscription",
		zap.String("user_id", req.UserID.String()),
		zap.String("service", req.ServiceName),
		zap.Int("price", req.Price),
		zap.String("start_date", req.StartDate))

	sub, err := toModel(req)
	if err != nil {
		s.logger.Warn("invalid subscription data",
			zap.String("user_id", req.UserID.String()),
			zap.Error(err))

		if errors.Is(err, ErrEndDateBeforeStart) {
			return SubscriptionResp{}, ErrEndDateBeforeStart
		}
		return SubscriptionResp{}, ErrInvalidDateFormat
	}

	created, err := s.repo.CreateSubscription(ctx, sub)
	if err != nil {
		return SubscriptionResp{}, err
	}
	response := toResponse(created)
	return response, nil
}

func (s *Service) GetSubscription(ctx context.Context, id int64) (SubscriptionResp, error) {
	s.logger.Info("getting subscription",
		zap.Int64("subscription_id", id))

	rawResponse, err := s.repo.GetSubscription(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			s.logger.Debug("subscription not found", zap.Int64("id", id))
			return SubscriptionResp{}, ErrSubscriptionNotFound
		}
		s.logger.Error("failed to get subscription from repo",
			zap.Int64("id", id),
			zap.Error(err))
		return SubscriptionResp{}, err
	}
	response := toResponse(rawResponse)

	return response, nil
}

func (s *Service) DeleteSubscription(ctx context.Context, id int64) error {
	s.logger.Info("deleting subscription",
		zap.Int64("subscription_id", id))

	if err := s.repo.DeleteSubscription(ctx, id); err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			s.logger.Debug("subscription not found", zap.Int64("id", id))
			return ErrSubscriptionNotFound
		}
		s.logger.Error("failed to delete subscription from repo",
			zap.Int64("id", id),
			zap.Error(err))
		return err
	}

	return nil
}

func toModel(req SubscriptionReq) (Subscription, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return Subscription{}, ErrInvalidDateFormat
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return Subscription{}, ErrInvalidDateFormat
		}
		endDate = &parsed
	}

	if endDate != nil && endDate.Before(startDate) {
		return Subscription{}, ErrEndDateBeforeStart
	}

	return Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func toResponse(rawSub Subscription) SubscriptionResp {
	startDate := rawSub.StartDate.Format("01-2006")

	var endDate *string
	if rawSub.EndDate != nil {
		formatted := rawSub.EndDate.Format("01-2006")
		endDate = &formatted
	}

	return SubscriptionResp{
		ID:          rawSub.ID,
		ServiceName: rawSub.ServiceName,
		Price:       rawSub.Price,
		UserID:      rawSub.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   rawSub.CreatedAt,
		UpdatedAt:   rawSub.UpdatedAt,
	}
}
