package subscription

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents the database model for subscriptions
type Subscription struct {
	ID          int64      `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name" example:"Yandex Plus"`
	Price       int        `db:"price" json:"price" example:"400"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `db:"start_date" json:"start_date" example:"07-2025"`
	EndDate     *time.Time `db:"end_date" json:"end_date" example:"12-2025"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at" example:"2025-02-19T10:00:00Z"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at" example:"2025-02-19T10:00:00Z"`
}

// SubscriptionReq represents the request body for creating/updating a subscription
type SubscriptionReq struct {
	ServiceName string    `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int       `json:"price" binding:"required" example:"400"`
	UserID      uuid.UUID `json:"user_id" binding:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     *string   `json:"end_date" example:"12-2026"`
}

// SubscriptionResp represents the response body for subscription data
type SubscriptionResp struct {
	ID          int64      `json:"id" example:"42"`
	ServiceName string     `json:"service_name" example:"Yandex Plus"`
	Price       int        `json:"price" example:"400"`
	UserID      uuid.UUID  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string     `json:"start_date" example:"07-2025"`
	EndDate     *string    `json:"end_date" example:"12-2026"`
	CreatedAt   *time.Time `json:"created_at" example:"2025-02-19T10:00:00Z"`
	UpdatedAt   *time.Time `json:"updated_at" example:"2025-02-19T10:00:00Z"`
}

// ListResult represents paginated list of subscriptions
type ListResult struct {
	Subscriptions []SubscriptionResp `json:"subscriptions"`
	Total         int                `json:"total" example:"1000"`
	Offset        int                `json:"offset" example:"10"`
	Limit         int                `json:"limit" example:"50"`
	HasMore       bool               `json:"has_more" example:"true"`
}

// TotalCostReq represents query parameters for total cost calculatio
type TotalCostReq struct {
	ServiceName *string `form:"service_name" example:"Yandex Plus"`
	UserID      *string `form:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `form:"start_date" binding:"required" example:"01-2025"`
	EndDate     string  `form:"end_date" binding:"required" example:"12-2025"`
}

// TotalCostResp represents the response for total cost calculation
type TotalCostResp struct {
	TotalCostReq
	TotalCost          int `json:"total_cost" example:"4800"`
	SubscriptionsCount int `json:"subscriptions_count" example:"4"`
}

// Filters represents the internal filters for database queries
// Used in total cost calculation to pass parsed and validated data to repository
type Filters struct {
	ServiceName *string    `json:"service_name" example:"Yandex Plus"`
	UserID      *uuid.UUID `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `json:"start_date" example:"07-2025"`
	EndDate     time.Time  `json:"end_date" example:"12-2025"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"invalid date format, use MM-YYYY"`
}
