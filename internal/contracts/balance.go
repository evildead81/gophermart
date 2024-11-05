package contracts

type Balance struct {
	CurrentBalance float64 `json:"current"`
	TotalWithdrawn float64 `json:"withdrawn"`
}
