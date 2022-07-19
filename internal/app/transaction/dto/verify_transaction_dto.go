package dto

import (
	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

type VerifyTransactionDto struct {
	ID         string             `json:"ID"`
	Session    domain.UserSession `json:"session"`
	Credential string             `json:"credential"`
}
