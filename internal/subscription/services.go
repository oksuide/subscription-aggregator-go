package subscription

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		zap.String("start_date", req.StartDate),
		zap.Stringp("end_date", req.EndDate))

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
		s.logger.Error("failed to create subscription in repo",
			zap.String("user_id", req.UserID.String()),
			zap.String("service", req.ServiceName),
			zap.Int("price", req.Price),
			zap.String("start_date", req.StartDate),
			zap.Stringp("end_date", req.EndDate),
			zap.Error(err))
		return SubscriptionResp{}, err
	}
	return toResponse(created), nil
}

func (s *Service) GetSubscription(ctx context.Context, id int64) (SubscriptionResp, error) {
	s.logger.Info("getting subscription",
		zap.Int64("subscription_id", id))

	subscription, err := s.repo.GetSubscription(ctx, id)
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
	return toResponse(subscription), nil
}

func (s *Service) UpdateSubscription(ctx context.Context, id int64, req SubscriptionReq) (SubscriptionResp, error) {
	s.logger.Info("updating subscription",
		zap.Int64("subscription_id", id),
		zap.String("service", req.ServiceName),
		zap.String("user_id", req.UserID.String()),
		zap.Int("price", req.Price),
		zap.String("start_date", req.StartDate),
		zap.Stringp("end_date", req.EndDate))

	updates, err := toModel(req)
	if err != nil {
		s.logger.Warn("invalid subscription data",
			zap.String("user_id", req.UserID.String()),
			zap.Error(err))

		if errors.Is(err, ErrEndDateBeforeStart) {
			return SubscriptionResp{}, ErrEndDateBeforeStart
		}
		return SubscriptionResp{}, ErrInvalidDateFormat
	}

	updated, err := s.repo.UpdateSubscription(ctx, id, updates)
	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			s.logger.Debug("subscription not found", zap.Int64("id", id))
			return SubscriptionResp{}, ErrSubscriptionNotFound
		}
		s.logger.Error("failed to update subscription in repo",
			zap.Int64("id", id),
			zap.Error(err))
		return SubscriptionResp{}, err
	}
	return toResponse(updated), nil
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

func (s *Service) ListSubscriptions(ctx context.Context, offset int, limit int) (ListResult, error) {
	s.logger.Info("listing subscriptions",
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	subscriptions, hasmore, total, err := s.repo.ListSubscriptions(ctx, offset, limit)
	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			s.logger.Debug("subscription not found",
				zap.Int("offset", offset),
				zap.Int("limit", limit))
			return ListResult{}, ErrSubscriptionNotFound
		}
		s.logger.Error("failed to list subscription from repo",
			zap.Int("offset", offset),
			zap.Int("limit", limit),
			zap.Error(err))
		return ListResult{}, err
	}
	return toResponseSlice(subscriptions, total, offset, limit, hasmore), nil
}

func (s *Service) CalculateTotalCost(ctx context.Context, req TotalCostReq) (TotalCostResp, error) {
	s.logger.Info("calculating total cost",
		zap.Stringp("service_name", req.ServiceName),
		zap.Stringp("user_id", req.UserID),
		zap.String("start_date", req.StartDate),
		zap.String("end_date", req.EndDate))

	filters, err := toFilters(req)
	if err != nil {
		s.logger.Warn("invalid subscription data",
			zap.Stringp("user_id", req.UserID),
			zap.Error(err))

		if errors.Is(err, ErrEndDateBeforeStart) {
			return TotalCostResp{}, ErrEndDateBeforeStart
		}
		return TotalCostResp{}, ErrInvalidDateFormat
	}

	cost, count, err := s.repo.GetTotalCost(ctx, filters)
	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			s.logger.Debug("no subscriptions found for period",
				zap.Stringp("service_name", req.ServiceName),
				zap.Stringp("user_id", req.UserID),
				zap.String("start_date", req.StartDate),
				zap.String("end_date", req.EndDate))
			return TotalCostResp{}, ErrSubscriptionNotFound
		}
		s.logger.Error("failed to calculate cost",
			zap.Stringp("service_name", req.ServiceName),
			zap.Stringp("user_id", req.UserID),
			zap.String("start_date", req.StartDate),
			zap.String("end_date", req.EndDate),
			zap.Error(err))
		return TotalCostResp{}, err
	}

	response := TotalCostResp{
		TotalCost:          cost,
		SubscriptionsCount: count,
		TotalCostReq: TotalCostReq{
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		},
	}

	if req.UserID != nil {
		response.UserID = req.UserID
	}
	if req.ServiceName != nil {
		response.ServiceName = req.ServiceName
	}

	s.logger.Info("cost calculated successfully",
		zap.Int("total_cost", cost),
		zap.Int("subscriptions_count", count))

	return response, nil
}

func toModel(req SubscriptionReq) (Subscription, error) {
	startDate, _, err := ParseDate(req.StartDate)
	if err != nil {
		return Subscription{}, err
	}
	_, endDate, err := ParseDate(req.EndDate)
	if err != nil {
		return Subscription{}, err
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

func toResponse(subscription Subscription) SubscriptionResp {
	startDate := subscription.StartDate.Format("01-2006")

	var endDate *string
	if subscription.EndDate != nil {
		formatted := subscription.EndDate.Format("01-2006")
		endDate = &formatted
	}

	return SubscriptionResp{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}
}

func toFilters(req TotalCostReq) (Filters, error) {
	startDate, _, err := ParseDate(req.StartDate)
	if err != nil {
		return Filters{}, err
	}

	endDate, _, err := ParseDate(req.EndDate)
	if err != nil {
		return Filters{}, err
	}

	if endDate.Before(startDate) {
		return Filters{}, ErrEndDateBeforeStart
	}

	var uuidID *uuid.UUID
	if req.UserID != nil {
		parsedID, err := uuid.Parse(*req.UserID)
		if err != nil {
			return Filters{}, fmt.Errorf("failed to parse user_id: %w", err)
		}
		uuidID = &parsedID
	}

	return Filters{
		ServiceName: req.ServiceName,
		UserID:      uuidID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func ParseDate(input any) (time.Time, *time.Time, error) {
	switch v := input.(type) {
	case string:
		parsed, err := time.Parse("01-2006", v)
		if err != nil {
			return time.Time{}, nil, ErrInvalidDateFormat
		}
		return parsed, nil, nil

	case *string:
		if v == nil {
			return time.Time{}, nil, nil
		}
		parsed, err := time.Parse("01-2006", *v)
		if err != nil {
			return time.Time{}, nil, ErrInvalidDateFormat
		}
		return time.Time{}, &parsed, nil

	default:
		return time.Time{}, nil, ErrInvalidDateFormat
	}
}

func toResponseSlice(subscriptions []Subscription, total, offset, limit int, hasMore bool) ListResult {
	result := make([]SubscriptionResp, len(subscriptions))
	for i, v := range subscriptions {
		result[i] = toResponse(v)
	}
	return ListResult{
		Subscriptions: result,
		Total:         total,
		Offset:        offset,
		Limit:         limit,
		HasMore:       hasMore,
	}
}
