package models

import (
	"fmt"

	"github.com/keslerliv/ilia-project/wallet/internal/entities"
	"github.com/keslerliv/ilia-project/wallet/pkg/db"
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

func PostTransaction(action string, value int64, uid int64) (*entities.Transaction, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Lock wallet to prevent race conditions
	var currentBalance int64
	err = tx.QueryRow(`SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE`, uid).Scan(&currentBalance)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("duplicate transaction detected")
	}

	// Calculate new balance
	var newBalance int64
	switch action {
	case "DEBIT":
		newBalance = currentBalance - value
		if newBalance < 0 {
			return nil, fmt.Errorf("insufficient balance")
		}
	case "CREDIT":
		newBalance = currentBalance + value
	default:
		return nil, fmt.Errorf("invalid type")
	}

	// Update wallet balance
	_, err = tx.Exec(`UPDATE wallets SET balance = $1 WHERE user_id = $2`, newBalance, uid)
	if err != nil {
		return nil, err
	}

	// Register transaction
	var transaction entities.Transaction
	err = tx.QueryRow(
		`INSERT INTO wallet_transactions (user_id, amount, action) VALUES ($1, $2, $3) RETURNING id, user_id, amount, action`,
		uid, value, action,
	).Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Type)
	if err != nil {
		return nil, err
	}

	// Commit if everything is successful
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func GetTransactions(uid int64) ([]entities.Transaction, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT id, user_id, amount, action FROM wallet_transactions WHERE user_id = $1", uid)
	var transactions []entities.Transaction

	for rows.Next() {
		var t entities.Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Type)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

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
