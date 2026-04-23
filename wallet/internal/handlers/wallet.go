package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/keslerliv/ilia-project/wallet/internal/entities"
	"github.com/keslerliv/ilia-project/wallet/internal/models"
)

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(entities.UIDKey).(int64)
	var payload entities.PostTransactionPayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transaction, err := models.PostTransaction(payload.Type, payload.Amout, user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(entities.UIDKey).(int64)

	transactions, err := models.GetTransactions(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

func GetBalance(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(entities.UIDKey).(int64)

	ballance, err := models.GetBalance(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{"amount": ballance}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
