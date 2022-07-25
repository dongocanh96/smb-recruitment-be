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
	pinCredentialManager domain.PinCredentialManager
}

func NewCreateTransactionService(userSession domain.UserSessionHelper, otpCredentialManager domain.OtpCredentialManager,
	pinCredentialManager domain.PinCredentialManager) CreateTransactionService {
	return &CreateTransactionServiceImp{userSession: userSession, otpCredentialManager: otpCredentialManager,
		pinCredentialManager: pinCredentialManager}
}

func (service *CreateTransactionServiceImp) Invoke(dto *dto.CreateTransactionDto, r context.Context) (string, error) {
	userSession, err := service.userSession.GetFromContext(r)
	if err != nil {
		return "", err
	}

	if err := service.validate(dto, userSession); err != nil {
		return "", err
	}

	transaction := &domain.Transaction{
		ID:                  uuid.New().String(),
		UserID:              userSession.ID,
		State:               domain.WaitAuthorization,
		AuthorizationMethod: alias.AuthMethods[dto.AuthMethod],
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
	if err := TransactionCodeValidate(dto.TransactionCode); err != nil {
		return err
	}

	if err := AmountValidate(dto.Amount); err != nil {
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

func TransactionCodeValidate(transactionCode string) error {
	if _, ok := alias.ValidTransactionCode[transactionCode]; ok {
		return nil
	}

	return alias.ErrMessageTransactionCodeNotFound
}

func AmountValidate(amount float64) error {
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
	if user.ConfiguredTransactionCredential.Otp != nil &&
		authMethod != alias.AuthMethod1 {
		return alias.ErrMessageMethodNotAllow
	}

	if user.ConfiguredTransactionCredential.Pin != nil &&
		authMethod != alias.AuthMethod2 {
		return alias.ErrMessageMethodNotAllow
	}

	if authMethod != alias.AuthMethod1 && authMethod != alias.AuthMethod2 {
		return alias.ErrMessageMethodNotAllow
	}

	return nil
}
func (service *CreateTransactionServiceImp) RequestNewOtp(userSession domain.UserSession) error {
	if !userSession.ConfiguredTransactionCredential.IsOtpConfigured() {
		return domain.ErrOtpNotConfigured
	}
	return nil
}
