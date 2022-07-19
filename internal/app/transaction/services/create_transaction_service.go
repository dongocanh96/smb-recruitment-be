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

	transaction := &domain.Transaction{
		ID:                 uuid.New().String(),
		UserID:             userSession.ID,
		State:              domain.Success,
		TransactionCode:    dto.TransactionCode,
		Amount:             dto.Amount,
		SourceAccount:      userSession.AccountReference,
		DestinationAccount: dto.DestinationAccount,
		CreatedAt:          time.Now().UTC(),
	}

	if dto.AuthMethod == alias.AuthMethod1 {
		transaction.AuthorizationMethod = domain.OtpAuthorization
	} else {
		transaction.AuthorizationMethod = domain.PinAuthorization
	}

	if err = pg.Wrap(nil).Save(transaction); err != nil {
		return "", err
	}

	return transaction.ID, nil
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
	for _, code := range alias.ValidTransactionCode {
		if transactionCode == code {
			return nil
		}
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
	for _, destinationCode := range alias.ValidDestination {
		if destination == destinationCode {
			return nil
		}
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

	return nil
}
