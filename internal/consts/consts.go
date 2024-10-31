package consts

type cookieName string
type contextKey string

const (
	CookieName string     = "AuthToken"
	UserIDKey  contextKey = "userID"
)
