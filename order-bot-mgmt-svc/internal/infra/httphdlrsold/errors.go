package httphdlrsold

import "net/http"

const (
	ErrMsgFailedMarshalResponse     = "failed to marshal response"
	ErrMsgFailedCheckDatabaseHealth = "failed to check database health"
	ErrMsgFailedMarshalHealthCheck  = "failed to marshal health check response"
	ErrMsgInvalidRequestBody        = "invalid request body"
	LogMsgFailedWriteResponse       = "failed to write response"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
