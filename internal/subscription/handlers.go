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
	rg.GET("/subscriptions", h.GetAllSubscriptions)
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
	// ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	// defer cancel()
	// response, err := h.service.CreateSubscription(ctx, req)
	// if err != nil {
	// }
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
func (h *Handler) GetAllSubscriptions(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	// defer cancel()
	// response, err := h.service.CreateSubscription(ctx, req)
	// if err != nil {
	// }
}

func getID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	return strconv.ParseInt(idStr, 10, 64)
}
