package models

import (
	"fmt"
	"strings"

	"github.com/keslerliv/user/internal/entities"
	"github.com/keslerliv/user/pkg/db"
)

func PostUser(user *entities.User) (*entities.ResponseUserData, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var userData entities.ResponseUserData
	err = conn.QueryRow(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, first_name, last_name, email",
		user.FirstName, user.LastName, user.Email, user.Password,
	).Scan(&userData.ID, &userData.FirstName, &userData.LastName, &userData.Email)
	if err != nil {
		return nil, err
	}

	return &userData, nil
}

func GetUsers() (users []entities.ResponseUserData, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	sql := `SELECT id, first_name, last_name, email FROM users`

	rows, err := conn.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user entities.ResponseUserData
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
		)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func GetUser(userID int64) (user entities.ResponseUserData, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	sql := `SELECT id, first_name, last_name, email FROM "users" WHERE id = $1`
	err = conn.QueryRow(sql, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		return
	}

	return
}

func PatchUser(userID int64, payload map[string]interface{}) (user entities.ResponseUserData, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	// Build dinamic query
	query := `UPDATE users SET `
	args := []interface{}{}
	i := 1

	if v, ok := payload["first_name"]; ok && v != "" {
		query += fmt.Sprintf("first_name = $%d,", i)
		args = append(args, v)
		i++
	}

	if v, ok := payload["last_name"]; ok && v != "" {
		query += fmt.Sprintf("last_name = $%d,", i)
		args = append(args, v)
		i++
	}

	if v, ok := payload["email"]; ok && v != "" {
		query += fmt.Sprintf("email = $%d,", i)
		args = append(args, v)
		i++
	}

	if v, ok := payload["password"]; ok && v != "" {
		// 🔐 aqui você deveria hashar antes
		query += fmt.Sprintf("password = $%d,", i)
		args = append(args, v)
		i++
	}
	query = strings.TrimSuffix(query, ",")

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, first_name, last_name, email", i)
	args = append(args, userID)

	// Run query
	err = conn.QueryRow(query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if err != nil {
		return
	}

	return
}

func DeleteUser(userID int64) (err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	result, err := conn.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return
}

func GetUserByEmail(email string) (user entities.User, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	sql := `SELECT id, first_name, last_name, email, password FROM "users" WHERE email = $1`
	err = conn.QueryRow(sql, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return
	}

	return
}
