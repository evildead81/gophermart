package storages

import "github.com/evildead81/gophermart/internal/contracts"

type Storage interface {
	RegisterUser(login string, password string) error
	CheckUserCredentials(login string, password string) error
	GetUserIDByLogin(login string) (int64, error)
	GetUserIDByOrderNumber(orderNumber string) (int64, error)
	CreateOrder(userID int, orderNumber string) error
	GetUserOrders(userID int) ([]contracts.Order, error)
	GetUserBalance(userID int) (contracts.Balance, error)
	Withdraw(userID int, order string, sum float64) error
	GetUserWithdrawals(userID int) ([]contracts.Withdrawal, error)
	ProcessAccruals(accrualAddress string)
}
