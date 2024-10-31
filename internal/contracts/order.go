package contracts

type Order struct {
	OrderNumber string  `json:"number"`
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual,omitempty"`
	UploadedAt  string  `json:"uploaded_at"`
}