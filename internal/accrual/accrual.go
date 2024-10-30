package accrual

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evildead81/gophermart/internal/contracts"
)

func GetAccrualStatus(orderNumber string, accrualSystemAddress string) (*contracts.AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", accrualSystemAddress, orderNumber)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Accrual service returned status: %v", resp.StatusCode)
	}

	var accrualResp contracts.AccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
		return nil, err
	}
	return &accrualResp, nil
}
