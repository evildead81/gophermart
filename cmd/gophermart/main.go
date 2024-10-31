package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/evildead81/gophermart/internal/config"
	"github.com/evildead81/gophermart/internal/handlers"
	"github.com/evildead81/gophermart/internal/middlewares"
	"github.com/evildead81/gophermart/internal/storages"
	dbstorage "github.com/evildead81/gophermart/internal/storages/db-storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func startAccrualProcessing(storage storages.Storage, accrualSystemAddress string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			storage.ProcessAccruals(accrualSystemAddress)
		}
	}
}

func main() {
	config, err := config.GetServerConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("pgx", config.DBUri)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := dbstorage.New(db)

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", handlers.Register(storage))
			r.Post("/login", handlers.Login(storage))
			r.Group(func(r chi.Router) {
				r.Use(middlewares.AuthMiddleware)
				r.Route("/orders", func(r chi.Router) {
					r.Get("/", handlers.GetOrders(storage))
					r.Post("/{number}", handlers.CreateOrder(storage))
				})
				r.Route("/balance", func(r chi.Router) {
					r.Get("/", handlers.GetBalance(storage))
					r.Post("/withdraw", handlers.Withdraw(storage))
				})
				r.Get("/withdrawals", handlers.GetWithdrawals(storage))
			})
		})
	})

	go startAccrualProcessing(storage, config.AccrualSystemAddress, 5*time.Minute)

	srv := &http.Server{
		Addr:    config.RunAddress,
		Handler: r,
	}
	srv.ListenAndServe()
}
