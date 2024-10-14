package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents the structure of a user
type User struct {
    ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Username string             `json:"username" bson:"username"`
    Password string             `json:"password" bson:"password"`
}

// Note represents the structure of a note
type Note struct {
    ID      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
    UserID  string             `json:"user_id" bson:"user_id"`
    Content string             `json:"content" bson:"content"`
}
