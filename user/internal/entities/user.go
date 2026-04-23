package entities

import "github.com/google/uuid"

// Main user model
type User struct {
	ID string `json:"id"`
}

func (a *User) NewUser(uid string) *User {
	return &User{
		ID: uuid.New().String(),
	}
}

// Request user model
type CreateUserPayload struct {
	UserID string `json:"user_id"`
}
