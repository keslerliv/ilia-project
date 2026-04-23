package entities

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UIDKey contextKey = "user_id"

type Claims struct {
	UserID json.Number `json:"user_id"`
	jwt.RegisteredClaims
}

// Main wallet model
type Wallet struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	Balance int `json:"balance"`
}

type Transaction struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Amount int64  `json:"amount"`
	Type   string `json:"type"`
}

// Request wallet model
type PostTransactionPayload struct {
	UserID int64  `json:"user_id"`
	Type   string `json:"type"`
	Amout  int64  `json:"amount"`
}
