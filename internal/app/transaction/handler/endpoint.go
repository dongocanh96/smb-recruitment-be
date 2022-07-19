package handler

import (
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
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
	requestDto := &dto.CreateTransactionDto{}

	if err := requestDto.Bind(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, &TransactionHandlerFailed{Message: err.Error()})
		return
	}

	transactionID, err := transactionEndpoint.transactionService.CreateTransaction(requestDto, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, &TransactionHandlerFailed{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)

	render.JSON(w, r, &CreateTransactionSuccess{transactionID})
}

func (transactionEndpoint *TransactionEndpoint) HandleVerifyTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userSession, err := transactionEndpoint.userSessionHelper.GetFromContext(r.Context())
	if err != nil {
		render.JSON(w, r, &TransactionHandlerFailed{Message: err.Error()})
	}

	verifyTransaction := VerifyTransactionRequest{}
	err = verifyTransaction.Bind(r)
	if err != nil {
		render.JSON(w, r, &TransactionHandlerFailed{Message: err.Error()})
	}

	transaction := &dto.VerifyTransactionDto{
		ID:         id,
		Session:    userSession,
		Credential: verifyTransaction.Credential,
	}

	err = transactionEndpoint.transactionService.VerifyTransaction(transaction, r.Context())
	if err != nil {
		render.JSON(w, r, &TransactionHandlerFailed{Message: err.Error()})
	}
	render.JSON(w, r, &VerifyTransactionSuccess{})
}

func (transactionEndpoint *TransactionEndpoint) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {

}
