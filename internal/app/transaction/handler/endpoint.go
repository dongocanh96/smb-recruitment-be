package handler

import (
	"github.com/tunaiku/mobilebanking/internal/app/transaction/services"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/pkg/jwt"
)

type TransactionEndpoint struct {
	userSessionHelper  domain.UserSessionHelper
	transactionService services.TransactionCompositionService
}

func NewTransactionEndpoint(
	userSessionHelper domain.UserSessionHelper,
	transactionCompositionService services.TransactionCompositionService) *TransactionEndpoint {
	return &TransactionEndpoint{
		userSessionHelper:  userSessionHelper,
		transactionService: transactionCompositionService,
	}
}

func (transactionEndpoint *TransactionEndpoint) BindRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r = jwt.WrapChiRouterWithAuthorization(r)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			})
		})
		r.Post("/transaction", transactionEndpoint.HandleCreateTransaction)
		r.Put("/transaction/{id}/verify", transactionEndpoint.HandleVerifyTransaction)
		r.Get("/transaction/{id}", transactionEndpoint.HandleGetTransaction)
	})
}

func (transactionEndpoint *TransactionEndpoint) HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	render.JSON(w, r, &CreateTransactionSuccess{})
}

func (transactionEndpoint *TransactionEndpoint) HandleVerifyTransaction(w http.ResponseWriter, r *http.Request) {

	render.JSON(w, r, &VerifyTransactionSuccess{})
}

func (transactionEndpoint *TransactionEndpoint) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {

}
