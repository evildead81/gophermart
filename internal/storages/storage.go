package storages

type Storage interface {
	IsUserExists(login string) bool
	CreateUser(login string, password string) error
}
