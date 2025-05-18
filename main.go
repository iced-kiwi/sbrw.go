package main

import (
	"fmt"
	"gosbrw/config"
	"gosbrw/database"
	"gosbrw/server"
	"gosbrw/server/routes"
	"log"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully.")
	log.Printf("Using HTTP Port from GameConfig: %d", config.AppConfig.Game.Port)

	dbCfg := config.AppConfig.Database
	if err := database.InitializeDatabase(dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName); err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	log.Println("PostgreSQL database initialized successfully.")
	defer database.CloseDB()

	engine := server.NewEngine()
	log.Println("Server engine created.")

	engine.RegisterRoutes(routes.ConfigureRoutes)
	log.Println("Routes registered.")

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.Game.IP, config.AppConfig.Game.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := engine.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
