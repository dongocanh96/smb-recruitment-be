package alias

import "errors"

const (
	ValidTransactionCode1 string  = "T001"
	MinimumAmount         float64 = 10000000
	ValidTransactionCode2 string  = "T002"
	ValidDestination1     string  = "10001"
	ValidDestination2     string  = "10002"
	AuthMethod2           string  = "PIN"
	AuthMethod1           string  = "OTP"
)

var (
	ErrMessageMethodNotAllow          = errors.New("auth method not allowed")
	ErrMessageTransactionCodeNotFound = errors.New("transaction code not found")
	ErrMessageAmountTooLow            = errors.New("amount does not reach the minimum transaction amount")
	ErrMessageDestinationNotFound     = errors.New("destination account not found")
)
