package main

import (
	"log"

	"github.com/Elvis-Benites-N/GolangChat/db"
	"github.com/Elvis-Benites-N/GolangChat/internal/user"
	"github.com/Elvis-Benites-N/GolangChat/internal/ws"
	"github.com/Elvis-Benites-N/GolangChat/router"
)

func main() {
	// Initialize the database connection
	dbConn, err := db.NewDatabaseLocal()
	if err != nil {
		log.Fatalf("Could not initialize the database connection: %s", err)
	}

	// Create the user repository
	userRep := user.NewRepository(dbConn.GetDBLocal())

	// Create the user service
	userSvc := user.NewService(userRep)

	// Create the user handler
	userHandler := user.NewHandler(userSvc)

	// Create the WebSocket hub
	hub := ws.NewHub()

	// Create the WebSocket handler
	wsHandler := ws.NewHandler(hub)

	// Run the WebSocket hub
	go hub.Run()

	// Initialize and start the router
	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:4200")
}
