package main

import (
	"log"
	"os"
	"todo/db"
	"todo/handlers"
	"todo/middlewares"
	"todo/models"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	db.ConnectMongo()
	db.ConnectPostgres()
	if err := db.ConnectRedis(); err != nil {
        panic(err)
    }

	if err := db.PG.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		log.Fatalf("migrate error: %v", err)
	}

	r := gin.Default()

	r.POST("register", handlers.Register)
	r.POST("login", handlers.Login)

	r.POST("/refresh", handlers.Refresh)
    r.POST("/logout", handlers.Logout)

	api := r.Group("/tasks", middlewares.Auth())


	{
		api.POST("/", handlers.CreateTask)
		api.GET("/", handlers.GetTasks)
		api.GET("/:id", handlers.GetTaskByID)
		api.PUT("/:id", handlers.UpdateTask)
		api.DELETE("/:id", handlers.DeleteTask)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("server on :" + port)
	r.Run(":" + port)



}