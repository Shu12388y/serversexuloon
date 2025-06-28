package pkg

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// MongoDBClientConnection returns a reference to the MongoDB database
func MongoDBClientConnection() *mongo.Database {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	var DB_URI = os.Getenv("DBURI")

	if DB_URI == "" {
		log.Fatal("❌ Database URI is required")
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Ping failed:", err)
	}

	fmt.Println("✅ Connected to MongoDB")

	// Return a reference to the database
	return client.Database("testdb")
}
