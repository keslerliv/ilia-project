package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/keslerliv/wallet/internal/entities"
	"github.com/keslerliv/wallet/internal/models"
)

func PostValue(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(entities.UIDKey).(int64)
	var payload entities.PostBallancePayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ballance, err := models.PostValue(payload.Action, payload.Value, user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{"new_ballance": float64(ballance) / 100}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetBalance(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(entities.UIDKey).(int64)

	ballance, err := models.GetBalance(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{"ballance": float64(ballance) / 100}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
