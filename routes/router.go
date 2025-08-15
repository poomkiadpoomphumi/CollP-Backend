package routes

import (
	controller "collp-backend/controllers"
	"collp-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	public := r.Group("/api")
	{
		// Google OAuth routes
		public.GET("/auth/google/login", gin.WrapF(controller.GoogleLogin))
		public.GET("/auth/google/callback", gin.WrapF(controller.GoogleCallback))

		// CollP auth routes
		public.POST("/collp/login", gin.WrapF(controller.CollPLogin))
		public.POST("/collp/register", gin.WrapF(controller.CollPRegister))
	}

	// Private routes (with authentication)
	private := r.Group("/api")
	private.Use(gin.WrapH(middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a placeholder handler, the actual handler will be called by Gin
	}))))
	{
		private.GET("/collp/main-menu", gin.WrapF(controller.MainMenu))
	}
}
