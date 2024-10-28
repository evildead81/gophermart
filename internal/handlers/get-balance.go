package handlers

import (
	"net/http"

	"github.com/evildead81/gophermart/internal/storages"
)

func GetBalance(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

	}
}
