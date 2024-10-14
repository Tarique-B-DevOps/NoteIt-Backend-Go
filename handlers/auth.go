package handlers

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "backend-golang/database"
)

type User struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    Username string `json:"username" bson:"username"`
    Password string `json:"password" bson:"password"`
}

// Register a new user
func Register(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existingUser User
    err := database.UserCol().FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existingUser)
    if err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
        return
    } else if err != mongo.ErrNoDocuments {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking username availability"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }
    user.Password = string(hashedPassword)

    _, err = database.UserCol().InsertOne(context.Background(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login a user
func Login(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var dbUser User
    err := database.UserCol().FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&dbUser)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": dbUser.ID})
}
