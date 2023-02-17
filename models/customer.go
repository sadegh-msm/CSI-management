package models

type Customer struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Credit   float64 `json:"credit"`
}
