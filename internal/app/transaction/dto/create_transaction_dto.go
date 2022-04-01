package dto

import "github.com/tunaiku/mobilebanking/internal/app/domain"

type CreateTransactionDto struct {
	Code        string
	Amount      float64
	Destination string
	AuthMethod  string
	Session     domain.UserSession
}
