package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"my-rest-api/config"
	"my-rest-api/models"
)

func GetUsers(c *gin.Context) {
	collection := config.Client.Database("test").Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching users"})
		return
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	for cursor.Next(context.TODO()) {
		var user models.User
		cursor.Decode(&user)
		// Don't send password in response
		user.Password = ""
		users = append(users, user)
	}

	c.JSON(200, users)
}
