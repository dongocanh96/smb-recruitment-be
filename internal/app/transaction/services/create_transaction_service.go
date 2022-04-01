package services

import (
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
	"github.com/tunaiku/mobilebanking/internal/pkg/pg"
)

type CreateTransactionService interface {
	Invoke(dto *dto.CreateTransactionDto) (string, error)
}

type CreateTransactionServiceImp struct {
}

func NewCreateTransactionService() CreateTransactionService {
	return &CreateTransactionServiceImp{}
}

func (service *CreateTransactionServiceImp) Invoke(dto *dto.CreateTransactionDto) (string, error) {
	if err := service.validate(dto); err != nil {
		return "", err
	}

	transaction := domain.Transaction{}

	return transaction.ID, pg.Wrap(nil).Save(&transaction)
}

func (service *CreateTransactionServiceImp) validate(dto *dto.CreateTransactionDto) error {

	return nil
}
