package entities

import "github.com/google/uuid"

// Main wallet model
type Wallet struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (a *Wallet) NewWallet(uid string) *Wallet {
	return &Wallet{
		ID:     uuid.New().String(),
		UserID: uid,
	}
}

// Request wallet model
type CreateWalletPayload struct {
	UserID string `json:"user_id"`
}
