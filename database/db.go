package database

import (
    "context"
    "os"
    "log"
    "time"
    "fmt"
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
    environment := os.Getenv("ENVIRONMENT")
    
    if environment == "Prod" {
        username := os.Getenv("MONGO_USERNAME")
        password := os.Getenv("MONGO_PASSWORD")
        cluster := os.Getenv("MONGO_CLUSTER_URI")
        
        mongoURI = fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=prod-cluster",
            username, password, cluster)
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
