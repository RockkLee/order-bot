package httphdlrsold

import (
	"net/http"
	"order-bot-mgmt-svc/internal/apperr"
)

var (
	ErrMsgFailedMarshalResponse = apperr.Err{
		Code: "ErrMsgFailedMarshalResponse",
		Msg:  "failed to marshal response",
	}
	ErrMsgFailedCheckDatabaseHealth = apperr.Err{
		Code: "ErrMsgFailedCheckDatabaseHealth",
		Msg:  "failed to check database health",
	}
	ErrMsgFailedMarshalHealthCheck = apperr.Err{
		Code: "ErrMsgFailedMarshalHealthCheck",
		Msg:  "failed to marshal health check response",
	}
	ErrMsgInvalidRequestBody = apperr.Err{
		Code: "ErrMsgInvalidRequestBody",
		Msg:  "invalid request body",
	}
)

const LogMsgFailedWriteResponse = "failed to write response"

func WriteError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
