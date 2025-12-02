package main

import (
	"HOSEROF_API/config"
	"HOSEROF_API/routes"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("warning: .env not found or failed to load")
	}

	config.InitFirebase()
	config.InitSupabase()
	defer config.DB.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("Server running on :", port)
	router := routes.SetupRouter()
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("server failed:", err)
	}
}
