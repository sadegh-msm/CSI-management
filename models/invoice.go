package models

import (
	"time"
)

type Invoice struct {
	ID             int64     `json:"id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Amount         float64   `json:"amount"`
	SubscriptionID int64     `json:"subscription_id"`
}
