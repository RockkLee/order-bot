package httphdlrs

import "net/http"

const (
	ErrMsgInvalidRequestBody        = "invalid request body"
	ErrMsgUserAlreadyExists         = "user already exists"
	ErrMsgInvalidCredentials        = "invalid credentials"
	ErrMsgFailedCreateUser          = "failed to create user"
	ErrMsgFailedLogin               = "failed to login"
	ErrMsgInvalidRefreshToken       = "invalid refresh token"
	ErrMsgLoggedOut                 = "logged out"
	ErrMsgFailedMarshalResponse     = "failed to marshal response"
	ErrMsgFailedCheckDatabaseHealth = "failed to check database health"
	ErrMsgFailedMarshalHealthCheck  = "failed to marshal health check response"
	LogMsgFailedWriteResponse       = "failed to write response"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
