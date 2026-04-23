package entities

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/keslerliv/user/config"
	"golang.org/x/crypto/bcrypt"
)

// Main user model
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *User) NewUser(email, password string) (*User, error) {
	encp, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(encp),
	}, nil
}

func (u *User) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)) == nil
}

func (u *User) CreateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": u.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	secret := config.Env.JWTSecret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Request user model
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
