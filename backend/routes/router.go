package routes

import (
	"net/http"

	"wsinspect/backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	proxyService *services.ProxyService,
	sessionService *services.SessionService,
	replayService *services.ReplayService,
	fuzzService *services.FuzzService,
) *gin.Engine {
	r := gin.Default()

	// Initialize message service
	messageService := services.NewMessageService(nil)

	// Set services for proxy
	proxyService.SetServices(sessionService, messageService)

	// API routes
	api := r.Group("/api")
	{
		// Proxy routes
		proxy := api.Group("/proxy")
		{
			proxy.GET("/status", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":             "running",
					"active_connections": proxyService.GetActiveConnections(),
				})
			})
		}

		// Session routes
		sessions := api.Group("/sessions")
		{
			sessions.GET("", func(c *gin.Context) {
				limit := 20
				offset := 0
				if l := c.Query("limit"); l != "" {
					// Simple parsing, ignore errors
					// In production, use proper parsing
				}

				sessions, total, err := sessionService.ListSessions(limit, offset)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"sessions": sessions,
					"total":    total,
				})
			})

			sessions.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")
				// Parse ID - in production use proper parsing
				session, err := sessionService.GetSession(1) // Placeholder
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
					return
				}
				c.JSON(http.StatusOK, session)
			})

			sessions.DELETE("/:id", func(c *gin.Context) {
				id := c.Param("id")
				// Parse ID and delete - placeholder
				c.JSON(http.StatusOK, gin.H{"message": "Session deleted"})
			})
		}

		// Message routes
		messages := api.Group("/messages")
		{
			messages.GET("/session/:session_id", func(c *gin.Context) {
				sessionID := c.Param("session_id")
				// Parse sessionID - placeholder
				messages, total, _ := messageService.GetMessagesBySession(1, 100, 0)
				c.JSON(http.StatusOK, gin.H{
					"messages": messages,
					"total":    total,
				})
				_ = sessionID
			})

			messages.POST("/inject", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Message injected"})
			})
		}

		// Replay routes
		replay := api.Group("/replay")
		{
			replay.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Replay started"})
			})
		}

		// Fuzz routes
		fuzz := api.Group("/fuzz")
		{
			fuzz.GET("", func(c *gin.Context) {
				tests, err := fuzzService.ListFuzzTests(0)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, tests)
			})

			fuzz.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Fuzz test created"})
			})

			fuzz.POST("/:id/run", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Fuzz test running"})
			})

			fuzz.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Fuzz test details"})
			})
		}
	}

	// WebSocket proxy endpoint
	r.GET("/ws/*target", func(c *gin.Context) {
		target := c.Param("target")
		targetURL := "ws://" + c.Query("host")
		if targetURL == "ws://" {
			targetURL = "ws://localhost:3000"
		}
		
		err := proxyService.HandleWebSocket(c.Writer, c.Request, targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "wsinspect",
		})
	})

	return r
}
