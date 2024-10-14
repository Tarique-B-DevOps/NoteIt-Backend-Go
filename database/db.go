package database

import (
    "context"
    "os"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client *mongo.Client
    ctx    = context.Background()
)

// UserCol will return the users collection
func UserCol() *mongo.Collection {
    return client.Database("noteit").Collection("users")
}

// NoteCol will return the notes collection
func NoteCol() *mongo.Collection {
    return client.Database("noteit").Collection("notes")
}

// ConnectDatabase connects to MongoDB
func ConnectDatabase() {
    var mongoURI string
    if os.Getenv("ENVIRONMENT") == "Prod" {
        mongoURI = os.Getenv("MONGO_URI")
    } else {
        mongoURI = "mongodb://localhost:27017"
    }

    var err error
    client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
}

// CheckConnection checks if the database is reachable
func CheckConnection() error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Set a timeout
    defer cancel()
    
    return client.Ping(ctx, nil) // Ping the MongoDB server
}

// DisconnectDatabase disconnects from MongoDB
func DisconnectDatabase() {
    if err := client.Disconnect(ctx); err != nil {
        log.Fatalf("Failed to disconnect from MongoDB: %v", err)
    }
}
