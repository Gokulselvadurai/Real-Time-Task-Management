package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"my-rest-api/config"
	"my-rest-api/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestBody struct {
	Tasks []struct {
		TaskName    string `json:"task_name"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
		AssignedTo  string `json:"assigned_to"`
	} `json:"tasks"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = primitive.NewObjectID()

	// Validate assignedTo user exists
	if task.AssignedTo != "" {
		collection := config.Client.Database("test").Collection("users")
		var user models.User
		err := collection.FindOne(context.TODO(), bson.M{"username": task.AssignedTo}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Assigned user not found"})
			return
		}
	}

	_, err := config.TaskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetTasks(c *gin.Context) {
	// Retrieve the email from the context
	email, _ := c.Get("email")

	// Retrieve the username associated with the email from the users collection
	var userResult bson.M
	err := config.Client.Database("test").Collection("users").
		FindOne(context.TODO(), bson.M{"email": email}).Decode(&userResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user by email"})
		return
	}

	// Extract the username from the user document
	// username, ok := userResult["username"].(string)

	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Username not found"})
	// 	return
	// }

	// Now query the tasks based on the email or username
	cursor, err := config.TaskCollection.Find(
		context.TODO(),
		bson.M{
			"userEmail": email, // Match by email only
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []models.Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func AskChatGPT(c *gin.Context) {
	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	// Prepare the task details for the prompt
	var taskDetails string
	for _, task := range body.Tasks {
		taskDetails += fmt.Sprintf(`
    Task: %s
    Description: %s
    Due Date: %s
    Assigned To: %s`, task.TaskName, task.Description, task.DueDate, task.AssignedTo)
	}

	// Create the full prompt
	prompt := fmt.Sprintf(`
    You are a helpful assistant. Here are the details of some tasks:
    %s
    Can you provide some advice on how to finish these tasks efficiently?`, taskDetails)

	// OpenAI request body
	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  150,
		"temperature": 0.7,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing the request"})
		return
	}

	// Retrieve the OpenAI API key from the environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "API Key not set"})
		return
	}

	// Send the POST request to OpenAI API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error communicating with OpenAI"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing OpenAI response"})
		return
	}

	var aiResponse OpenAIResponse
	if err := json.Unmarshal(bodyBytes, &aiResponse); err != nil {
		log.Println("Error unmarshalling response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing OpenAI response"})
		return
	}

	// Check if there are any choices in the response
	if len(aiResponse.Choices) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No response from OpenAI"})
		return
	}

	// Return the first choice from OpenAI response
	c.JSON(http.StatusOK, gin.H{
		"message": aiResponse.Choices[0].Message.Content,
	})
}

func UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	err = config.TaskCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Toggle the status
	newStatus := !task.Status
	_, err = config.TaskCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": newStatus}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": taskID, "status": newStatus})
}
