package services

import (
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
)

type VerifyTransactionService interface {
	Invoke(dto *dto.VerifyTransactionDto) error
}

type VerifyTransactionServiceImp struct {
}

func NewVerifyTransactionService() VerifyTransactionService {
	return &VerifyTransactionServiceImp{}
}

func (service *VerifyTransactionServiceImp) Invoke(dto *dto.VerifyTransactionDto) error {

	return nil
}
