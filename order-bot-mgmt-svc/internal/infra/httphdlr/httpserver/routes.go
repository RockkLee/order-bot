package httpserver

import (
	"fmt"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/httphdlr"
	"time"

	"github.com/gin-gonic/gin"
)

func Run(s *Server, ginMode string, addr string) {
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(ginMode)
	routers := gin.New()
	routers.Use(gin.Recovery())
	routers.Use(corsMiddleware())
	routers.Use(gin.LoggerWithFormatter(func(p gin.LogFormatterParams) string {
		if p.StatusCode < 400 {
			return ""
		}
		return fmt.Sprintf("[%s] %d %s %s ip=%s latency=%s err=%s\n",
			p.TimeStamp.Format(time.RFC3339Nano), p.StatusCode, p.Method, p.Path, p.ClientIP, p.Latency, p.ErrorMessage,
		)
	}))

	routers.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})
	routers.GET("/health", func(c *gin.Context) {
		stats, err := s.dbService().Health()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": httphdlr.ErrMsgFailedCheckDatabaseHealth})
			return
		}
		c.JSON(http.StatusOK, stats)
	})

	root := routers.Group("/orderbotmgmt")
	public := root.Group("")
	protected := root.Group("")
	protected.Use(authMiddleware(s))
	auth := public.Group(httphdlr.AuthPrefix)
	httphdlr.RegisterAuthRoutes(auth, s)
	menus := protected.Group(httphdlr.MenuPrefix)
	httphdlr.RegisterMenuRoutes(menus, s)
	bot := protected.Group(httphdlr.BotPrefix)
	httphdlr.RegisterBotRoutes(bot, s)

	err := routers.Run(addr)
	if err != nil {
		panic(err.Error())
	}
}
