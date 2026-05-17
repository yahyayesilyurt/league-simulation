package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yahyayesilyurt/league-simulation/config"
	"github.com/yahyayesilyurt/league-simulation/internal/handler"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    config.ConnectDatabase()
    config.SeedDatabase(config.DB)
    
    r := handler.SetupRouter(config.DB)

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    r.Run(":" + port)
}