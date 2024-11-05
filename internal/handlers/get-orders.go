package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evildead81/gophermart/internal/consts"
	"github.com/evildead81/gophermart/internal/storages"
)

func GetOrders(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(consts.UserIDKey).(int)

		orders, err := storage.GetUserOrders(userID)
		if err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		json.NewEncoder(rw).Encode(orders)
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
	}
}
