package services

import (
	"context"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
)

type VerifyTransactionService interface {
	Invoke(dto *dto.VerifyTransactionDto, r context.Context) error
}

type VerifyTransactionServiceImp struct {
	isOtp domain.OtpCredentialManager
	isPin domain.PinCredentialManager
}

func NewVerifyTransactionService(isOtp domain.OtpCredentialManager, isPin domain.PinCredentialManager) VerifyTransactionService {
	return &VerifyTransactionServiceImp{isOtp: isOtp, isPin: isPin}
}

func (service *VerifyTransactionServiceImp) Invoke(dto *dto.VerifyTransactionDto, r context.Context) error {

	return nil
}
