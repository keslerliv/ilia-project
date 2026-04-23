package entities

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/keslerliv/ilia-project/user/config"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UIDKey contextKey = "user_id"

type Claims struct {
	UserID json.Number `json:"user_id"`
	jwt.RegisteredClaims
}

// Main user model
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (a *User) NewUser(user CreateUserPayload) (*User, error) {
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("email and password are required")
	}

	if user.FirstName == "" || user.LastName == "" {
		return nil, errors.New("first name and last name are required")
	}

	encp, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  string(encp),
	}, nil
}

func (a *User) NewPatchUser(payload CreateUserPayload) (*User, error) {
	user := &User{}

	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Email = payload.Email

	if payload.Password != "" {
		encp, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
		if err != nil {
			return nil, err
		}
		user.Password = string(encp)
	}

	return user, nil
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
	User UserLoginPayload `json:"user"`
}

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserPayload struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type ResponseUserData struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}
