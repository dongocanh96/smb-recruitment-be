package handler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/micro/go-micro/v3/errors"
	"math/big"
	"net/http"
)

type TransactionHandlerFailed struct {
	HttpCode int    `json:"-"`
	Message  string `json:"message"`
	Error    error  `json:"-"`
}

func (resp *TransactionHandlerFailed) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(resp.HttpCode)
	return nil
}

func (resp *TransactionHandlerFailed) RenderWithError(w http.ResponseWriter, r *http.Request) {
	switch v := resp.Error.(type) {
	case *errors.Error:
		_ = render.Render(w, r, &TransactionHandlerFailed{Message: v.Detail, HttpCode: int(v.Code)})
	default:
		_ = render.Render(w, r, &TransactionHandlerFailed{Message: v.Error(), HttpCode: http.StatusInternalServerError})
	}
}

type CreateTransactionRequest struct {
	TransactionCode     string     `json:"transaction_code"`
	Amount              *big.Float `json:"amount"`
	DestinationAccount  string     `json:"destination_account"`
	AuthorizationMethod string     `json:"auth_method"`
}

func (payload *CreateTransactionRequest) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(payload); err != nil {
		return err
	}
	return nil
}

type CreateTransactionSuccess struct {
	TransactionID string `json:"transaction_id"`
}

func (resp *CreateTransactionSuccess) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusCreated)
	return nil
}

type VerifyTransactionRequest struct {
	Credential string `json:"credential"`
}

func (payload *VerifyTransactionRequest) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(payload); err != nil {
		return err
	}
	return nil
}

type VerifyTransactionSuccess struct {
	TransactionID string `json:"transaction_id"`
}

func (resp *VerifyTransactionSuccess) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusAccepted)
	return nil
}

type GetTransactionSuccess struct {
	ID                 string  `json:"id"`
	Amount             float64 `json:"amount"`
	DestinationAccount string  `json:"destination_account"`
	State              string  `json:"state"`
}

func (resp *GetTransactionSuccess) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}
