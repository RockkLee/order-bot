package httpserver

import (
	"fmt"
	"log/slog"
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
	orders := protected.Group(httphdlr.OrderPrefix)
	httphdlr.RegisterOrderRoutes(orders, s)

	health := public.Group("/health")
	health.GET("/chk", func(c *gin.Context) {
		stats, err := s.dbService().Health()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": httphdlr.ErrMsgFailedCheckDatabaseHealth})
			return
		}
		slog.Info("httpserver.routes.Run.health(), stats: ", stats)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	err := routers.Run(addr)
	if err != nil {
		panic(err.Error())
	}
}
