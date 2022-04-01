package dto

import "github.com/tunaiku/mobilebanking/internal/app/domain"

type VerifyTransactionDto struct {
	ID 	string
	Session     domain.UserSession
	Credential string
}
