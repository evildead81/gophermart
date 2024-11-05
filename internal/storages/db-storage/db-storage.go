package dbstorage

import (
	"database/sql"
	"log"

	"github.com/evildead81/gophermart/internal/accrual"
	"github.com/evildead81/gophermart/internal/contracts"
	"github.com/evildead81/gophermart/internal/errors"
	"github.com/evildead81/gophermart/internal/hashing"
	"github.com/evildead81/gophermart/internal/storages"
)

type DBStorage struct {
	db      *sql.DB
	storage storages.Storage
}

func New(db *sql.DB) *DBStorage {
	storage := &DBStorage{
		db: db,
	}

	storage.initDB()
	return storage
}

func (s DBStorage) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
	    user_id SERIAL PRIMARY KEY,
	    login VARCHAR(255) UNIQUE NOT NULL,
	    password_hash TEXT NOT NULL,
	    created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS orders (
	    order_id SERIAL PRIMARY KEY,
	    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
	    order_number VARCHAR(50) UNIQUE NOT NULL,
	    status VARCHAR(20) CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')) DEFAULT 'NEW',
	    accrual NUMERIC(10, 2) DEFAULT 0,
	    uploaded_at TIMESTAMP DEFAULT NOW(),
	    processed_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS balances (
	    user_id INT PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
	    current_balance NUMERIC(10, 2) DEFAULT 0,
	    total_withdrawn NUMERIC(10, 2) DEFAULT 0,
	    last_updated TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS withdrawals (
	    withdrawal_id SERIAL PRIMARY KEY,
	    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
	    order_number VARCHAR(50) NOT NULL,
	    sum NUMERIC(10, 2) NOT NULL,
	    processed_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) CheckUserCredentials(login string, password string) error {
	var storedHash string
	err := s.db.QueryRow("SELECT password_hash FROM users WHERE login = $1", login).Scan(&storedHash)

	if err == sql.ErrNoRows {
		return errors.ErrInvalidCredentials
	} else if err != nil {
		return err
	}

	if !hashing.CheckPasswordHash(password, storedHash) {
		return errors.ErrInvalidCredentials
	}

	return nil
}

func (s *DBStorage) RegisterUser(login string, password string) error {
	var userExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)", login).Scan(&userExists)
	if err != nil {
		return err
	}

	if userExists {
		return errors.ErrUserIsAlreadyExists
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return err
	}

	var userID int
	err = tx.QueryRow("INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING user_id", login, hashedPassword).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO balances (user_id, current_balance, total_withdrawn) VALUES ($1, 0, 0)", userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) GetUserIDByLogin(login string) (int64, error) {
	var userID int
	err := s.db.QueryRow("SELECT user_id FROM users WHERE login = $1", login).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return int64(userID), nil
}

func (s *DBStorage) GetUserIDByOrderNumber(orderNumber string) (int64, error) {
	var userID int
	err := s.db.QueryRow("SELECT user_id FROM orders WHERE order_number = $1", orderNumber).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return int64(userID), nil
}

func (s *DBStorage) CreateOrder(userID int, orderNumber string) error {
	_, err := s.db.Exec("INSERT INTO orders (user_id, order_number, status) VALUES ($1, $2, 'NEW')", userID, orderNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) GetUserOrders(userID int) ([]contracts.Order, error) {
	rows, err := s.db.Query("SELECT order_number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()
	if err == sql.ErrNoRows {
		return make([]contracts.Order, 0), nil
	}

	var orders []contracts.Order
	for rows.Next() {
		var order contracts.Order
		if err := rows.Scan(&order.OrderNumber, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *DBStorage) GetUserBalance(userID int) (contracts.Balance, error) {
	var balance contracts.Balance
	err := s.db.QueryRow("SELECT current_balance, total_withdrawn FROM balances WHERE user_id = $1", userID).Scan(&balance.CurrentBalance, &balance.TotalWithdrawn)
	if err != nil {
		return contracts.Balance{}, err
	}

	return balance, nil
}

func (s *DBStorage) Withdraw(userID int, order string, sum float64) error {
	var currentBalance float64
	err := s.db.QueryRow("SELECT current_balance FROM balances WHERE user_id = $1", userID).Scan(&currentBalance)
	if err == sql.ErrNoRows || currentBalance < sum {
		return errors.ErrPaymentRequiredError
	} else if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE balances SET current_balance = current_balance - $1, total_withdrawn = total_withdrawn + $1 WHERE user_id = $2", sum, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)", userID, order, sum)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (s *DBStorage) GetUserWithdrawals(userID int) ([]contracts.Withdrawal, error) {
	rows, err := s.db.Query("SELECT order_number, sum, processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []contracts.Withdrawal

	for rows.Next() {
		var withdrawal struct {
			Order       string
			Sum         float64
			ProcessedAt string
		}
		if err := rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, contracts.Withdrawal{
			Order:       withdrawal.Order,
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (s *DBStorage) ProcessAccruals(accrualAddress string) {
	rows, err := s.db.Query("SELECT order_id, order_number, user_id FROM orders WHERE status IN ('NEW', 'PROCESSING')")
	if err != nil {
		log.Printf("Error fetching orders: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var orderID int
		var orderNumber string
		var userID int

		if err := rows.Scan(&orderID, &orderNumber, &userID); err != nil {
			log.Printf("Error scanning order: %v", err)
			continue
		}

		accrualResp, err := accrual.GetAccrualStatus(orderNumber, accrualAddress)
		if err != nil {
			log.Printf("Error getting accrual status: %v", err)
			continue
		}

		switch accrualResp.Status {
		case "PROCESSED":
			_, err := s.db.Exec(
				"UPDATE orders SET status = $1, accrual = $2, processed_at = NOW() WHERE order_id = $3",
				accrualResp.Status, accrualResp.Accrual, orderID,
			)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				continue
			}

			_, err = s.db.Exec(
				"UPDATE balances SET current_balance = current_balance + $1, last_updated = NOW() WHERE user_id = $2",
				accrualResp.Accrual, userID,
			)
			if err != nil {
				log.Printf("Error updating balance: %v", err)
			}

		case "INVALID":
			_, err := s.db.Exec(
				"UPDATE orders SET status = $1, processed_at = NOW() WHERE order_id = $2",
				accrualResp.Status, orderID,
			)
			if err != nil {
				log.Printf("Error updating invalid order status: %v", err)
			}

		case "PROCESSING", "REGISTERED":
			_, err := s.db.Exec(
				"UPDATE orders SET status = $1 WHERE order_id = $2",
				accrualResp.Status, orderID,
			)
			if err != nil {
				log.Printf("Error updating processing order status: %v", err)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return
	}
}
