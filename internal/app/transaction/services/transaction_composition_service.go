package services

import (
	"context"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/dto"
	"github.com/tunaiku/mobilebanking/internal/pkg/pg"
)

type TransactionCompositionService interface {
	CreateTransaction(dto *dto.CreateTransactionDto, ctx context.Context) (string, error)
	VerifyTransaction(dto *dto.VerifyTransactionDto, ctx context.Context) error
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

func (inst *TransactionCompositionServiceImp) CreateTransaction(dto *dto.CreateTransactionDto, ctx context.Context) (string, error) {
	return inst.createTransactionService.Invoke(dto, ctx)
}

func (inst *TransactionCompositionServiceImp) VerifyTransaction(dto *dto.VerifyTransactionDto, ctx context.Context) error {
	return inst.verifyTransactionService.Invoke(dto, ctx)
}

func (inst *TransactionCompositionServiceImp) GetTransaction(id string) (domain.Transaction, error) {
	transaction := domain.Transaction{ID: id}
	err := pg.Wrap(nil).Load(&transaction)
	return transaction, err
}
