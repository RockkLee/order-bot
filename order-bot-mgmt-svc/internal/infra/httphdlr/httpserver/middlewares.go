package httpserver

import (
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/httphdlr"
	"order-bot-mgmt-svc/internal/util/errutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	}
}

func authMiddleware(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, httphdlr.AuthPrefix+"/") || c.Request.URL.Path == httphdlr.AuthPrefix || c.Request.URL.Path == "/health" || c.Request.URL.Path == "/" {
			c.Next()
			return
		}

		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authService := s.AuthService()
		if authService == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err := authService.ValidateAccessToken(c.Request.Context(), accessToken); err != nil {
			slog.Debug(errutil.FormatErrChain(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func bearerToken(authHeader string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", false
	}
	t := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if t == "" {
		return "", false
	}
	return t, true
}
