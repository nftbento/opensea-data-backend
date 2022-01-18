/*

 */

package server

import (
	"github.com/NFTActions/opensea-data-backend/utils/ratelimit"
	"github.com/gin-gonic/gin"
)

func NewRouter(server *Server) *gin.Engine {
	gin.SetMode("debug")
	r := gin.Default()

	r.Use(ratelimit.GinMiddleware())
	r.Use(CORSMiddleware())
	r.GET("/ping", server.controller.base.HandlePing)

	v1 := r.Group("/v1")
	WithAdminRoutes(v1, server)
	WithUnauthorizedRoutes(v1, server)

	return r
}

func WithAdminRoutes(r *gin.RouterGroup, server *Server) {
	//todo: config admin routes based on config
	admin := r.Group("/admin")
	admin.Use(adminAuth())

	admin.POST("/activity/recent", server.controller.acti.HandleActivityCreate)
}

func WithUnauthorizedRoutes(r *gin.RouterGroup, server *Server) {
	r.GET("/activity/summary", server.controller.acti.HandleGetActivitySummary)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Auth-Token, Authorization, Code, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT , PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func adminAuth() gin.HandlerFunc {
	accounts := gin.Accounts{
		"larry":     "larrykey",
		"scheduler": "schedulerkey",
	}
	return gin.BasicAuth(accounts)
}
