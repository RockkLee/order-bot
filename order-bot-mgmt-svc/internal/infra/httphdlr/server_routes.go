package httphdlr

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouters(s *ServerContainer, ginMode string) http.Handler {
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
	auth := public.Group(AuthPrefix)
	RegisterAuthRoutes(auth, s)
	menus := protected.Group(MenuPrefix)
	RegisterMenuRoutes(menus, s)
	bot := protected.Group(BotPrefix)
	RegisterBotRoutes(bot, s)
	orders := protected.Group(OrderPrefix)
	RegisterOrderRoutes(orders, s)

	health := public.Group("/health")
	health.GET("/chk", func(c *gin.Context) {
		stats, err := s.db().Health()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": ErrMsgFailedCheckDatabaseHealth})
			return
		}
		slog.Debug("httpserver.routes.Run.health()", "stats", stats)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return routers
}
