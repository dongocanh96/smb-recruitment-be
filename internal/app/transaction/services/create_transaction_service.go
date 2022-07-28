package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/alias"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
	"github.com/tunaiku/mobilebanking/internal/pkg/pg"
)

type CreateTransactionService interface {
	Invoke(dto *dto.CreateTransactionDto, ctx context.Context) (string, error)
}

type CreateTransactionServiceImp struct {
	userSession          domain.UserSessionHelper
	otpCredentialManager domain.OtpCredentialManager
}

func NewCreateTransactionService(userSession domain.UserSessionHelper, otpCredentialManager domain.OtpCredentialManager,
) CreateTransactionService {
	return &CreateTransactionServiceImp{userSession: userSession, otpCredentialManager: otpCredentialManager}
}

func (service *CreateTransactionServiceImp) Invoke(dto *dto.CreateTransactionDto, r context.Context) (string, error) {
	userSession, err := service.userSession.GetFromContext(r)
	if err != nil {
		return "", err
	}

	if err := service.validate(dto, userSession); err != nil {
		return "", err
	}

	authMethod := domain.UnknownAuthorizationMethod
	if dto.AuthMethod == alias.AuthMethod1 {
		authMethod = domain.OtpAuthorization
	} else {
		authMethod = domain.PinAuthorization
	}

	transaction := &domain.Transaction{
		ID:                  uuid.New().String(),
		UserID:              userSession.ID,
		State:               domain.WaitAuthorization,
		AuthorizationMethod: authMethod,
		TransactionCode:     dto.TransactionCode,
		Amount:              dto.Amount,
		SourceAccount:       userSession.AccountReference,
		DestinationAccount:  dto.DestinationAccount,
		CreatedAt:           time.Now().UTC(),
	}

	if err = pg.Wrap(nil).Save(transaction); err != nil {
		return "", err
	}

	return transaction.ID, nil
}

func (service *CreateTransactionServiceImp) validate(dto *dto.CreateTransactionDto, userSession domain.UserSession) error {
	if err := ValidateTransactionCode(dto.TransactionCode); err != nil {
		return err
	}

	if err := ValidateAmount(dto.Amount); err != nil {
		return err
	}

	if err := CheckDestination(dto.DestinationAccount); err != nil {
		return err
	}

	if err := CheckValidMethod(dto.AuthMethod, userSession); err != nil {
		return err
	}

	if err := service.RequestNewOtp(userSession); err != nil && dto.AuthMethod == alias.AuthMethod1 {
		return err
	}
	return nil
}

func ValidateTransactionCode(transactionCode string) error {
	if _, ok := alias.ValidTransactionCode[transactionCode]; ok {
		return nil
	}

	return alias.ErrMessageTransactionCodeNotFound
}

func ValidateAmount(amount float64) error {
	if amount < alias.MinimumAmount {
		return alias.ErrMessageAmountTooLow
	}
	return nil
}

func CheckDestination(destination string) error {
	if _, ok := alias.ValidDestination[destination]; ok {
		return nil
	}

	return alias.ErrMessageDestinationNotFound
}

func CheckValidMethod(authMethod string, user domain.UserSession) error {
	if authMethod == alias.AuthMethod1 && user.ConfiguredTransactionCredential.Otp == nil {
		return alias.ErrMessageMethodNotConfigured
	}

	if authMethod == alias.AuthMethod2 && user.ConfiguredTransactionCredential.Pin == nil {
		return alias.ErrMessageMethodNotConfigured
	}

	if authMethod != alias.AuthMethod1 && authMethod != alias.AuthMethod2 {
		return alias.ErrMessageMethodNotSupported
	}

	return nil
}
func (service *CreateTransactionServiceImp) RequestNewOtp(userSession domain.UserSession) error {
	if !userSession.ConfiguredTransactionCredential.IsOtpConfigured() {
		return alias.ErrMessageOtpNotConfigured
	}
	return nil
}
