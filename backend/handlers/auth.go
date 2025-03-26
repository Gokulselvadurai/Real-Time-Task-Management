package handlers

import (
	"context"
	"fmt"
	"my-rest-api/config"
	"my-rest-api/models"
	"my-rest-api/utils"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}

	// Validate email format
	if !validateEmail(user.Email) {
		c.JSON(400, gin.H{"error": "Invalid email format"})
		return
	}

	// Check password length
	if len(user.Password) < 8 {
		c.JSON(400, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	// Check username length
	if len(user.Username) < 3 {
		c.JSON(400, gin.H{"error": "Username must be at least 3 characters long"})
		return
	}

	collection := config.Client.Database("test").Collection("users")

	// Check if email already exists
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(400, gin.H{"error": "Email already registered"})
		return
	}

	// Check if username already exists
	err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(400, gin.H{"error": "Username already taken"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user object with hashed password
	newUser := models.User{
		Email:    user.Email,
		Password: string(hashedPassword),
		Username: user.Username,
	}

	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error storing user"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	utils.SetTokenCookie(c, token)

	// Return user data without password
	c.JSON(200, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// Create a separate struct for signin requests
type SigninRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signin(c *gin.Context) {
	var signinReq SigninRequest

	if err := c.ShouldBindJSON(&signinReq); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}

	if !validateEmail(signinReq.Email) {
		c.JSON(400, gin.H{"error": "Invalid email format"})
		return
	}

	collection := config.Client.Database("test").Collection("users")
	var storedUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": signinReq.Email}).Decode(&storedUser)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(signinReq.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(signinReq.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	utils.SetTokenCookie(c, token)
	c.JSON(200, gin.H{
		"message": "Signin successful",
		"user": gin.H{
			"email":    storedUser.Email,
			"username": storedUser.Username,
		},
	})
}

func Signout(c *gin.Context) {
	utils.ClearTokenCookie(c)
	c.JSON(200, gin.H{"message": "Signout successful"})
}

func Protected(c *gin.Context) {
	email, _ := c.Get("email")

	// Get user data from database
	collection := config.Client.Database("test").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching user data"})
		return
	}

	c.JSON(200, gin.H{
		"message":  "This is a protected route",
		"email":    user.Email,
		"username": user.Username,
	})
}
