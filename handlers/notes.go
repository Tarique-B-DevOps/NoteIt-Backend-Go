package handlers

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "backend-golang/database"
)

type Note struct {
    ID      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
    UserID  string             `json:"user_id" bson:"user_id"`
    Content string             `json:"content" bson:"content"`
}

// Create a new note
func CreateNote(c *gin.Context) {
    var note Note
    if err := c.ShouldBindJSON(&note); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := database.NoteCol().InsertOne(context.Background(), note)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting note"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Note created successfully"})
}

// Get all notes for a user
func GetNotes(c *gin.Context) {
    userID := c.Param("user_id")

    cursor, err := database.NoteCol().Find(context.Background(), bson.M{"user_id": userID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving notes"})
        return
    }
    defer cursor.Close(context.Background())

    var notes []Note
    if err = cursor.All(context.Background(), &notes); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding notes"})
        return
    }

    c.JSON(http.StatusOK, notes)
}

// Delete a note
func DeleteNote(c *gin.Context) {
    noteID := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(noteID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
        return
    }

    _, err = database.NoteCol().DeleteOne(context.Background(), bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting note"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
