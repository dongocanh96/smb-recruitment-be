package alias

import (
	"errors"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

const (
	MinimumAmount     float64 = 3000
	AuthMethod1       string  = "otp"
	AuthMethod2       string  = "pin"
	TransactionCode1  string  = "T001"
	TransactionCode2  string  = "T002"
	Destination1      string  = "10001"
	Destination2      string  = "10002"
	DefaultOtp        string  = "111111"
	DefaultPin        string  = "123456"
	WaitAuthorization string  = "WaitAuthorization"
	Failed            string  = "Failed"
	Success           string  = "Success"
)

var AuthMethods = map[string]domain.AuthorizationMethod{
	AuthMethod1: domain.OtpAuthorization,
	AuthMethod2: domain.PinAuthorization,
}

var ValidTransactionCode = map[string]string{
	TransactionCode1: TransactionCode1,
	TransactionCode2: TransactionCode2,
}

var ValidDestination = map[string]string{
	Destination1: Destination1,
	Destination2: Destination2,
}

var TransactionState = map[domain.TransactionState]string{
	domain.WaitAuthorization: WaitAuthorization,
	domain.Success:           Success,
	domain.Failed:            Failed,
}

var (
	ErrMessageMethodNotConfigured     = errors.New("authorization method not configured")
	ErrMessageMethodNotSupported      = errors.New("unsupported authorization method")
	ErrMessageTransactionCodeNotFound = errors.New("transaction code not found")
	ErrMessageAmountTooLow            = errors.New("amount does not reach the minimum transaction amount")
	ErrMessageDestinationNotFound     = errors.New("destination account not found")
	ErrMessageOtpNotConfigured        = errors.New("OTP not configured")
	ErrMessagePinNotConfigured        = errors.New("PIN not configured")
	ErrMessageInvalidCredential       = errors.New("invalid credential")
	ErrMessageTransactionHadVerified  = errors.New("verification process already happened")
)
