package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "ok",
            "message": "League Simulation is running",
        })
    })

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    r.Run(":" + port)
}