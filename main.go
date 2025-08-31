package main

import (
	"github.com/gin-gonic/gin"
	"github.com/samansahebi/dekamond/controllers"
	"github.com/samansahebi/dekamond/models"
)

func init() {
	models.ConnectDatabase()
	models.DB.AutoMigrate(&models.User{}, &models.OTP{})
}

func main() {

	r := gin.Default()

	r.POST("/send-otp", controllers.SendOTP)
	r.POST("/verify-otp", controllers.VerifyOTP)

	protected := r.Group("/api")
	protected.Use(controllers.AuthMiddleware())

	protected.GET("/users", controllers.Users)
	protected.GET("/users/:id", controllers.GetUser)
	protected.GET("users/search", controllers.Search)

	r.Run(":8000")
}
