package dbstorage

import (
	"database/sql"

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

func (s *DBStorage) IsUserExists(login string) bool {
	var userExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)", login).Scan(&userExists)
	if err != nil {
		return false
	}
	return userExists
}

func (s *DBStorage) CreateUser(login string, password string) error {
	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO users (login, password_hash) VALUES ($1, $2)", login, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
