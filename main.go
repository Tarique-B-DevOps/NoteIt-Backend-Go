// main.go
package main

import (
    "backend-golang/database"
    "backend-golang/handlers"
    "net/http" // Import net/http for the health check
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    // Connect to MongoDB
    database.ConnectDatabase()
    defer database.DisconnectDatabase()

    // Set up the router
    router := gin.Default()

    // CORS configuration
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:          12 * time.Hour,
    }))

    // API endpoints
    router.POST("/login", handlers.Login)
    router.POST("/register", handlers.Register)
    router.POST("/notes", handlers.CreateNote)
    router.GET("/notes/:user_id", handlers.GetNotes)
    router.DELETE("/notes/:id", handlers.DeleteNote)

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "OK"})
    })

    // Database health check endpoint
    router.GET("/health/db", func(c *gin.Context) {
        err := database.CheckConnection()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "DB not reachable", "error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "DB is reachable"})
    })

    // Start the server
    router.Run(":8080")
}
