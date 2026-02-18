package subscription

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const requestTimeout = 5 * time.Second

type Handler struct {
	service *Service
	logger  *zap.Logger
}

func RegisterRoutes(rg *gin.RouterGroup, service *Service, logger *zap.Logger) {
	h := &Handler{service: service, logger: logger}
	rg.POST("/subscriptions", h.CreateSubscription)
	rg.GET("/subscriptions/:id", h.GetSubscription)
	rg.PUT("/subscriptions/:id", h.UpdateSubscription)
	rg.DELETE("/subscriptions/:id", h.DeleteSubscription)
	rg.GET("/subscriptions", h.ListSubscriptions)
	rg.GET("/subscriptions/total-cost", h.GetTotalCost)
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	var req SubscriptionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body",
			zap.Int("status", 400),
			zap.String("error", err.Error()),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method))
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid request body"})
		return
	}

	h.logger.Debug("creating subscription",
		zap.String("user_id", req.UserID.String()),
		zap.String("service", req.ServiceName),
		zap.Int("price", req.Price))

	response, err := h.service.CreateSubscription(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(504, gin.H{"error": "timeout"})
		case errors.Is(err, ErrInvalidDateFormat), errors.Is(err, ErrEndDateBeforeStart):
			c.JSON(400, gin.H{"error": err.Error()})
		default:
			h.logger.Error("failed to create subscription",
				zap.Int("status", 500),
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path))
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}
	h.logger.Info("subscription created successfully",
		zap.String("user_id", req.UserID.String()),
		zap.String("service", req.ServiceName),
		zap.Int("price", req.Price))

	c.JSON(201, response)
}

func (h *Handler) GetSubscription(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	id, err := getID(c)
	if err != nil {
		h.logger.Warn("invalid id format",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid subscription id"})
		return
	}

	h.logger.Debug("getting subscription",
		zap.Int64("subscription_id", id))

	response, err := h.service.GetSubscription(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(504, gin.H{"error": "timeout"})
		case errors.Is(err, ErrSubscriptionNotFound):
			c.JSON(404, gin.H{"error": "subscription not found"})
		default:
			h.logger.Error("failed to get subscription",
				zap.Int("status", 500),
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path))
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}
	h.logger.Info("subscription retrieved successfully",
		zap.Int64("id", id))

	c.JSON(200, response)
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	id, err := getID(c)
	if err != nil {
		h.logger.Warn("invalid id format",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid subscription id"})
		return
	}

	var req SubscriptionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body",
			zap.Int("status", 400),
			zap.String("error", err.Error()),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method))
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid request body"})
		return
	}

	h.logger.Debug("updating subscription",
		zap.Int64("subscription_id", id))

	response, err := h.service.UpdateSubscription(ctx, id, req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(504, gin.H{"error": "timeout"})
		case errors.Is(err, ErrInvalidDateFormat), errors.Is(err, ErrEndDateBeforeStart):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, ErrSubscriptionNotFound):
			c.JSON(404, gin.H{"error": "subscription not found"})
		default:
			h.logger.Error("failed to update subscription",
				zap.Int("status", 500),
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path))
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	h.logger.Info("subscription updated successfully",
		zap.Int64("id", id))

	c.JSON(200, response)
}

func (h *Handler) DeleteSubscription(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	id, err := getID(c)
	if err != nil {
		h.logger.Warn("invalid id format",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid subscription id"})
		return
	}
	h.logger.Debug("deleting subscription",
		zap.Int64("subscription_id", id))

	err = h.service.DeleteSubscription(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(504, gin.H{"error": "timeout"})
		case errors.Is(err, ErrSubscriptionNotFound):
			c.JSON(404, gin.H{"error": "subscription not found"})
		default:
			h.logger.Error("failed to delete subscription",
				zap.Int("status", 500),
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path))
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}
	h.logger.Info("subscription deleted successfully",
		zap.Int64("id", id))

	c.Status(204)
}

func (h *Handler) ListSubscriptions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		h.logger.Warn("invalid offset",
			zap.String("offset", c.Query("offset")),
			zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid offset"})
		return
	}

	if offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.logger.Warn("invalid limit",
			zap.String("limit", c.Query("limit")),
			zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid limit"})
		return
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	h.logger.Debug("listing subscriptions",
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	response, err := h.service.ListSubscriptions(ctx, offset, limit)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(504, gin.H{"error": "timeout"})
		case errors.Is(err, ErrSubscriptionNotFound):
			c.JSON(404, gin.H{"error": "subscriptions not found"})
		default:
			h.logger.Error("failed to list subscriptions",
				zap.Int("status", 500),
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path))
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	h.logger.Info("subscriptions listed successfully",
		zap.Int("count", response.Total),
		zap.Bool("has more", response.HasMore),
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	c.JSON(200, response)
}

func (h *Handler) GetTotalCost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	var req TotalCostReq
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Warn("invalid query parameters", zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid query parameters"})
		return
	}

	response, err := h.service.CalculateTotalCost(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(504, gin.H{"error": "timeout"})
		case errors.Is(err, ErrSubscriptionNotFound):
			c.JSON(404, gin.H{"error": "subscriptions not found"})
		case errors.Is(err, ErrInvalidDateFormat), errors.Is(err, ErrEndDateBeforeStart):
			c.JSON(400, gin.H{"error": err.Error()})
		default:
			h.logger.Error("failed to calculate cost",
				zap.Int("status", 500),
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path))
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	h.logger.Info("cost retrieved successfully",
		zap.Int("total_cost", response.TotalCost),
		zap.Int("subscriptions_count", response.SubscriptionsCount))

	c.JSON(200, response)
}

func getID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	return strconv.ParseInt(idStr, 10, 64)
}
