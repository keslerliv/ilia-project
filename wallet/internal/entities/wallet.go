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

// Request wallet model
type PostBallancePayload struct {
	Action string `json:"action"`
	Value  int64  `json:"value"`
}
