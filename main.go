package main

import (
	"log"
	"status-page-monitor/internal/database"
	"status-page-monitor/internal/env"
	"status-page-monitor/server/router"
)

func main() {
	env.InitEnv()

	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	router.InitRouter()

}
