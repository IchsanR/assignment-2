package main

import (
	"assigntment2/database"
	"assigntment2/routers"
)

func main() {
	var PORT = ":8080"
	database.DatabaseConnect()

	routers.StartServer().Run(PORT)
}
