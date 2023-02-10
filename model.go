package main

import "time"

// Transaction - transaction model fo each transaction
type Transaction struct {
	Amount    float64   `json:"amount"`
	TimeStamp time.Time `json:"timestamp"`
}

// City - city model to store the city name of end user
type City struct {
	City string `json:"city"`
}

// AllTransactions - model to store user details and transactions
type AllTransactions struct {
	UserId string `json:"user_id"`
	City
	Transactions      []Transaction `json:"transactions"`
	TotalAmount       float64
	TotalTransactions int64
	MinAmount         float64
	MaxAmount         float64
}

// Response model of statics api
type Statics struct {
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
	Count int64   `json:"count"`
}
