package httpserver

import (
	"net/http"
	"order-bot-mgmt-svc/internal/infra/httphdlr"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(), authMiddleware(s))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})
	r.GET("/health", func(c *gin.Context) {
		stats, err := s.dbService().Health()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": httphdlr.ErrMsgFailedCheckDatabaseHealth})
			return
		}
		c.JSON(http.StatusOK, stats)
	})

	auth := r.Group(httphdlr.AuthPrefix)
	httphdlr.RegisterAuthRoutes(auth, s)
	menus := r.Group(httphdlr.MenuPrefix)
	httphdlr.RegisterMenuRoutes(menus, s)
	bot := r.Group(httphdlr.BotPrefix)
	httphdlr.RegisterBotRoutes(bot, s)

	return r
}
