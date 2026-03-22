package api

import (
	"github.com/enterprise/enterprise-3tier/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires HTTP routes and middleware placeholders.
func RegisterRoutes(r *gin.Engine, user *handler.UserHandler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/users")
	{
		api.POST("", user.Create)
		api.GET("", user.List)
		api.GET("/:id", user.Get)
		api.PUT("/:id", user.Update)
		api.DELETE("/:id", user.Delete)
	}
}
