package router

import (
	"my-rest-api/handlers"
	"my-rest-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(middleware.CORS())

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/signin", handlers.Signin)
		auth.POST("/signout", handlers.Signout)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.Auth())
	{
		protected.GET("/protected", handlers.Protected)
		protected.POST("/tasks", handlers.CreateTask)
		protected.GET("/tasks", handlers.GetTasks)
		protected.PATCH("/tasks/:id/status", handlers.UpdateTaskStatus)
		protected.GET("/users", handlers.GetUsers)
		protected.POST("/ask-chatgpt", handlers.AskChatGPT)
	}

	return r
}
