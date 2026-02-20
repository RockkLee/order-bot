package httphdlr

import "order-bot-mgmt-svc/internal/apperr"

var (
	ErrMsgFailedCheckDatabaseHealth = apperr.Err{
		Code: "ErrMsgFailedCheckDatabaseHealth",
		Msg:  "failed to check database health",
	}
	ErrMsgInvalidRequestBody = apperr.Err{
		Code: "ErrMsgInvalidRequestBody",
		Msg:  "invalid request body",
	}
)
