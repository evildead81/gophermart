package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evildead81/gophermart/internal/consts"
	"github.com/evildead81/gophermart/internal/errors"
	"github.com/evildead81/gophermart/internal/storages"
)

func GetBalance(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(consts.UserIDKey).(int)
		balance, err := storage.GetUserBalance(userID)
		if err == errors.ErrNotFound {
			http.Error(rw, "Balance not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(rw, "Server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(balance)
		rw.WriteHeader(http.StatusOK)
	}
}
