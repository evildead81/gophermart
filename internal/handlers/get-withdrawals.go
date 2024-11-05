package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evildead81/gophermart/internal/consts"
	"github.com/evildead81/gophermart/internal/storages"
)

func GetWithdrawals(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(consts.UserIDKey).(int)

		withdrawals, err := storage.GetUserWithdrawals(userID)

		if err != nil {
			http.Error(rw, "Server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(withdrawals)
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
	}
}
