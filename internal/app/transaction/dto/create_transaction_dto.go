package dto

import (
	"encoding/json"
	"net/http"
)

type CreateTransactionDto struct {
	TransactionCode    string  `json:"transaction_code"`
	Amount             float64 `json:"amount"`
	DestinationAccount string  `json:"destination_account"`
	AuthMethod         string  `json:"auth_method"`
}

func (dto *CreateTransactionDto) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(dto); err != nil {
		return err
	}
	return nil
}
