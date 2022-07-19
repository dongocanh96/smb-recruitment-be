package alias

import "errors"

const (
	MinimumAmount float64 = 10000000
	AuthMethod1   string  = "OTP"
	AuthMethod2   string  = "PIN"
)

var ValidTransactionCode = map[int]string{
	1: "T001",
	2: "T002",
}

var ValidDestination = map[int]string{
	1: "10001",
	2: "10002",
}

var (
	ErrMessageMethodNotAllow          = errors.New("auth method not allowed")
	ErrMessageTransactionCodeNotFound = errors.New("transaction code not found")
	ErrMessageAmountTooLow            = errors.New("amount does not reach the minimum transaction amount")
	ErrMessageDestinationNotFound     = errors.New("destination account not found")
)
