package models

import (
	"fmt"

	"github.com/keslerliv/wallet/pkg/db"
)

func CreateWallet(uid int) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int64
	err = conn.QueryRow("INSERT INTO wallets (user_id) VALUES ($1) RETURNING id", uid).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func PostValue(action string, value int64, uid int64) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Lock wallet to prevent race conditions
	var currentBalance int64
	err = tx.QueryRow(`SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE`, uid).Scan(&currentBalance)
	if err != nil {
		return 0, err
	}

	// Check duplicate transaction in last 5 minutes
	var exists bool
	err = tx.QueryRow(`SELECT EXISTS (
			SELECT 1 FROM wallet_transactions
			WHERE user_id = $1
			  AND amount = $2
			  AND action = $3
			  AND created_at >= NOW() - INTERVAL '5 minutes')
	`, uid, value, action).Scan(&exists)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("duplicate transaction detected")
	}

	// Calculate new balance
	var newBalance int64
	switch action {
	case "withdraw":
		newBalance = currentBalance - value
		if newBalance < 0 {
			return 0, fmt.Errorf("insufficient balance")
		}
	case "deposit":
		newBalance = currentBalance + value
	default:
		return 0, fmt.Errorf("invalid action")
	}

	// Update wallet balance
	_, err = tx.Exec(`UPDATE wallets SET balance = $1 WHERE user_id = $2`, newBalance, uid)
	if err != nil {
		return 0, err
	}

	// Register transaction
	_, err = tx.Exec(`INSERT INTO wallet_transactions (user_id, amount, action) VALUES ($1, $2, $3)`, uid, value, action)
	if err != nil {
		return 0, err
	}

	// Commit if everything is successful
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newBalance, nil
}

// func PostValue(action string, value int, uid int64) (int64, error) {
// 	conn, err := db.OpenConnection()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer conn.Close()

// 	var balance int64

// 	if action == "withdraw" {
// 		err = conn.QueryRow("UPDATE wallets SET balance = balance - $1 WHERE user_id = $2 RETURNING balance", value, uid).Scan(&balance)
// 	} else if action == "deposit" {
// 		err = conn.QueryRow("UPDATE wallets SET balance = balance + $1 WHERE user_id = $2 RETURNING balance", value, uid).Scan(&balance)
// 	} else {
// 		return 0, fmt.Errorf("invalid action: %s", action)
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 		return 0, err
// 	}

// 	return balance, nil
// }

func GetBalance(uid int64) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var balance int64
	err = conn.QueryRow("SELECT balance FROM wallets WHERE user_id = $1", uid).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
