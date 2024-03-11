package main

import (
	"assigntment2/database"
	"assigntment2/routers"
	"log"
)

func main() {
	var PORT = ":8080"
	db, err := database.DatabaseConnect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	routers.StartServer(db).Run(PORT)
}
