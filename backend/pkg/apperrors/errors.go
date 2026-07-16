package apperrors

import "net/http"

type AppError struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Msg
}

func New(code int, message string) *AppError {
	return &AppError{
		Code: code,
		Msg:  message,
	}
}

var (
	ErrUserNotFound             = New(http.StatusNotFound, "user could not be found")
	ErrUserAlreadyExists        = New(http.StatusBadRequest, "user already exists")
	ErrInvalidCredentials       = New(http.StatusUnauthorized, "invalid credentials")
	ErrInsufficientBalance      = New(http.StatusBadRequest, "insufficient balance")
	ErrUserNotHaveWallet        = New(http.StatusBadRequest, "user does not have a wallet")
	ErrInvalidPassword          = New(http.StatusBadRequest, "invalid password")
	ErrWalletAlreadyExists      = New(http.StatusBadRequest, "wallet already exists")
	ErrTransactionNotFound      = New(http.StatusNotFound, "transaction not found")
	ErrTransactionStatusInvalid = New(http.StatusBadRequest, "transaction status is invalid")
	ErrInternalServer           = New(http.StatusInternalServerError, "internal server error")
)
