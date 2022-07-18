package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/alias"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
	"github.com/tunaiku/mobilebanking/internal/pkg/pg"
	"time"
)

type CreateTransactionService interface {
	Invoke(dto *dto.CreateTransactionDto, ctx context.Context) (string, error)
}

type CreateTransactionServiceImp struct {
	userSession domain.UserSessionHelper
	isOtp       domain.OtpCredentialManager
	isPin       domain.PinCredentialManager
}

func NewCreateTransactionService(userSession domain.UserSessionHelper, isOtp domain.OtpCredentialManager, isPin domain.PinCredentialManager) CreateTransactionService {
	return &CreateTransactionServiceImp{userSession: userSession, isOtp: isOtp, isPin: isPin}
}

func (service *CreateTransactionServiceImp) Invoke(dto *dto.CreateTransactionDto, r context.Context) (string, error) {
	userSession, err := service.userSession.GetFromContext(r)
	if err != nil {
		return "", err
	}

	if err := service.validate(dto, userSession); err != nil {
		return "", err
	}

	transaction := domain.Transaction{
		ID:                  uuid.New().String(),
		UserID:              userSession.ID,
		State:               1,
		AuthorizationMethod: 1,
		TransactionCode:     dto.TransactionCode,
		Amount:              dto.Amount,
		SourceAccount:       userSession.AccountReference,
		DestinationAccount:  dto.DestinationAccount,
		CreatedAt:           time.Now(),
	}

	return transaction.ID, pg.Wrap(nil).Save(&transaction)
}

func (service *CreateTransactionServiceImp) validate(dto *dto.CreateTransactionDto, userSession domain.UserSession) error {
	if err := GetTransactionPrivileges(dto.TransactionCode); err != nil {
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
	return nil
}

func GetTransactionPrivileges(transactionCode string) error {
	if transactionCode != alias.ValidTransactionCode1 && transactionCode != alias.ValidTransactionCode2 {
		return alias.ErrMessageTransactionCodeNotFound
	}
	return nil
}

func AmountValidate(amount float64) error {
	if amount < alias.MinimumAmount {
		return alias.ErrMessageAmountTooLow
	}
	return nil
}

func CheckDestination(destination string) error {
	if destination != alias.ValidDestination1 && destination != alias.ValidDestination2 {
		return alias.ErrMessageDestinationNotFound
	}
	return nil
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

	return nil
}
