package pkg

import (
	"log"
	"os"
)

func MongoDBClientConnection() {

	var DB_URI = os.Getenv("DBURI")

	if DB_URI == "" {
		log.Fatal("Database URI is required")
	}

}
