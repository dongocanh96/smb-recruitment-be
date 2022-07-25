package transaction

import (
	"github.com/tunaiku/mobilebanking/internal/app/transaction/services"
	"log"

	"github.com/go-chi/chi"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/handler"
	"go.uber.org/dig"
)

func Register(container *dig.Container) {
	container.Provide(func(userSession domain.UserSessionHelper, otpCredentialManager domain.OtpCredentialManager,
		pinCredentialManager domain.PinCredentialManager) services.CreateTransactionService {
		return services.NewCreateTransactionService(userSession, otpCredentialManager, pinCredentialManager)
	})

	container.Provide(func(userSession domain.UserSessionHelper, otpCredentialManager domain.OtpCredentialManager,
		pinCredentialManager domain.PinCredentialManager) services.VerifyTransactionService {
		return services.NewVerifyTransactionService(userSession, otpCredentialManager, pinCredentialManager)
	})

	container.Provide(func(
		createTransactionService services.CreateTransactionService,
		verifyTransactionService services.VerifyTransactionService) services.TransactionCompositionService {
		return services.NewTransactionCompositionService(createTransactionService, verifyTransactionService)
	})

	container.Provide(func(
		userSessionHelper domain.UserSessionHelper,
		transactionService services.TransactionCompositionService) *handler.TransactionEndpoint {
		return handler.NewTransactionEndpoint(userSessionHelper, transactionService)
	})
}

func Invoke(container *dig.Container) {
	err := container.Invoke(func(router chi.Router, endpoint *handler.TransactionEndpoint) {
		log.Println("invoke transaction startup ...")
		endpoint.BindRoutes(router)
	})
	if err != nil {
		log.Fatal(err)
	}
}
