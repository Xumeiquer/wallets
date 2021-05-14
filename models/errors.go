package models

const (
	MSG_STATUS_SUCCESS Errno = 1

	ERR_INSERT_WALLET Errno = 101
	ERR_NEW_WALLET    Errno = 102

	ERR_NOT_FOUND        Errno = 404
	ERR_DB_NOT_FOUND     Errno = 405
	ERR_WALLET_NOT_FOUND Errno = 406
	ERR_NOT_A_WALLET     Errno = 407
	ERR_UPDATE_WALLET    Errno = 408
	ERR_INSI_NOT_FOUND   Errno = 409
)

var (
	ERR_MESSAGES = map[Errno]string{
		MSG_STATUS_SUCCESS:   "operation successful",
		ERR_INSERT_WALLET:    "error saving the wallet",
		ERR_NEW_WALLET:       "invalid name",
		ERR_NOT_FOUND:        "element not found",
		ERR_DB_NOT_FOUND:     "database not avaliable in context",
		ERR_WALLET_NOT_FOUND: "wallet not found",
		ERR_NOT_A_WALLET:     "not a wallet",
		ERR_UPDATE_WALLET:    "unable to update the wallet",
		ERR_INSI_NOT_FOUND:   "inis not found",
	}
)

type Errno int

type ErrorMessage struct {
	Errno Errno
	Msg   string
}

func NewError(code Errno, msg string) ErrorMessage {
	return ErrorMessage{
		Errno: code,
		Msg:   msg,
	}
}
