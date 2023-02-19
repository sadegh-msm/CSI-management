package models

import "time"

type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Credit    float64   `json:"credit"`
	CreatedAt time.Time `json:"created_at"`
}
