package services

import (
	"context"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/alias"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
)

type VerifyTransactionService interface {
	Invoke(dto *dto.VerifyTransactionDto, r context.Context) error
}

type VerifyTransactionServiceImp struct {
	userSession          domain.UserSessionHelper
	otpCredentialManager domain.OtpCredentialManager
	pinCredentialManager domain.PinCredentialManager
}

func NewVerifyTransactionService(userSession domain.UserSessionHelper, otpCredentialManager domain.OtpCredentialManager,
	pinCredentialManager domain.PinCredentialManager) VerifyTransactionService {
	return &VerifyTransactionServiceImp{userSession: userSession, otpCredentialManager: otpCredentialManager,
		pinCredentialManager: pinCredentialManager}
}

func (service *VerifyTransactionServiceImp) Invoke(dto *dto.VerifyTransactionDto, r context.Context) error {
	userSession, err := service.userSession.GetFromContext(r)
	if err != nil {
		return err
	}

	if userSession.ConfiguredTransactionCredential.Otp != nil {
		return validateOtp(userSession, dto.Credential)
	}

	if userSession.ConfiguredTransactionCredential.Pin != nil {
		return validatePin(userSession, dto.Credential)
	}

	return nil
}

func validateOtp(userSession domain.UserSession, credential string) error {
	if !userSession.ConfiguredTransactionCredential.IsOtpConfigured() {
		return domain.ErrOtpNotConfigured
	}

	if credential != alias.DefaultOtp {
		return domain.ErrCredentialNotMatch
	}

	return nil

}

func validatePin(userSession domain.UserSession, credential string) error {
	if !userSession.ConfiguredTransactionCredential.IsPinConfigured() {
		return domain.ErrPinNotConfigured
	}

	if credential != alias.DefaultPin {
		return domain.ErrCredentialNotMatch
	}

	return nil
}
