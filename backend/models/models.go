package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Email    string `json:"email" binding:"required" bson:"email,unique"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required" bson:"username,unique"`
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	DueDate     time.Time          `json:"dueDate"`
	AssignedTo  string             `bson:"assignedto" json:"assignedto"`
	UserEmail   string             `bson:"userEmail" json:"userEmail"`
	Status      bool               `bson:"status" json:"status"` // Ensure status is included in the task list
}
