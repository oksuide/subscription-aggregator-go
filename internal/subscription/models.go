package subscription

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64      `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name"`
	Price       int        `db:"price" json:"price"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	StartDate   time.Time  `db:"start_date" json:"start_date"`
	EndDate     *time.Time `db:"end_date" json:"end_date"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
}

type SubscriptionReq struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date"`
}

type SubscriptionResp struct {
	ID          int64      `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   string     `json:"start_date"`
	EndDate     *string    `json:"end_date"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ListResult struct {
	Subscriptions []SubscriptionResp `json:"subscriptions"`
	Total         int                `json:"total"`
	Offset        int                `json:"offset"`
	Limit         int                `json:"limit"`
	HasMore       bool               `json:"has_more"`
}

type TotalCostReq struct {
	ServiceName *string `form:"service_name"`
	UserID      *string `form:"user_id"`
	StartDate   string  `form:"start_date" binding:"required"`
	EndDate     string  `form:"end_date" binding:"required"`
}

type Filters struct {
	ServiceName *string    `json:"service_name"`
	UserID      *uuid.UUID `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
}

type TotalCostResp struct {
	TotalCostReq
	TotalCost          int `json:"total_cost"`
	SubscriptionsCount int `json:"subscriptions_count"`
}
