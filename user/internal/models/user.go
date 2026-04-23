package models

import (
	"github.com/keslerliv/user/internal/entities"
	"github.com/keslerliv/user/pkg/db"
)

func CreateUser(user *entities.User) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int64
	err = conn.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetUserByEmail(email string) (user entities.User, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	sql := `SELECT id, email, password FROM "users" WHERE email = $1`
	err = conn.QueryRow(sql, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return
	}

	return
}
