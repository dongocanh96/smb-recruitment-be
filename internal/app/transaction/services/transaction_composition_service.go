package services

import (
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
)

type TransactionCompositionService interface {
	CreateTransaction(dto *dto.CreateTransactionDto) (string, error)
	VerifyTransaction(dto *dto.VerifyTransactionDto) error
	GetTransaction(id string) (domain.Transaction, error)
}

type TransactionCompositionServiceImp struct {
	createTransactionService CreateTransactionService
	verifyTransactionService VerifyTransactionService
}

func NewTransactionCompositionService(
	createTransactionService CreateTransactionService,
	verifyTransactionService VerifyTransactionService) TransactionCompositionService {
	return &TransactionCompositionServiceImp{
		createTransactionService: createTransactionService,
		verifyTransactionService: verifyTransactionService,
	}
}

func (inst *TransactionCompositionServiceImp) CreateTransaction(dto *dto.CreateTransactionDto) (string, error) {
	return inst.createTransactionService.Invoke(dto)
}

func (inst *TransactionCompositionServiceImp) VerifyTransaction(dto *dto.VerifyTransactionDto) error {
	return inst.verifyTransactionService.Invoke(dto)
}

func (inst *TransactionCompositionServiceImp) GetTransaction(id string) (domain.Transaction, error) {
	transaction := domain.Transaction{ID: id}

	return transaction, nil
}
