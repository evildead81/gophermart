package contracts

type OrderRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
