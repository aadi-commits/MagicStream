package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/database"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/routes"
)

func init() {

	//Load environment variable
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
}

func main(){

	//connect MongoDB 
	database.Connect()

	//Initialize controller dependencies
	controller.InitMovieController()
	controller.InitUserController()

	//Create gin router
	router:= gin.Default()

	router.GET("/hello", func(c *gin.Context){
		c.String(200, "Hello, Magic Stream Movies!")
	})

	api := router.Group("/api/v1")
	routes.SetupUnProtectedRoutes(api)
	routes.SetupProctectedRoutes(api)

	//Load API key
	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		fmt.Println("NVIDIA_API_KEY not set")
		os.Exit(1)
	}

	router.POST("/ai/infer", func(c *gin.Context) {

		var input struct {
			Prompt string `json:"prompt"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// 🔥 Correct OpenAI-compatible payload
		payload := map[string]interface{}{
			"model": "nvidia/nemotron-3-nano-30b-a3b",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": input.Prompt,
				},
			},
			"temperature": 1,
			"top_p":       1,
			"max_tokens":  1024,
		}

		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build request"})
			return
		}

		// ✅ Correct endpoint
		url := "https://integrate.api.nvidia.com/v1/chat/completions"

		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "AI service error",
				"details": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{
				"error":   "NVIDIA API returned error",
				"details": string(responseData),
			})
			return
		}

		c.Data(http.StatusOK, "application/json", responseData)
	})


	if err:= router.Run(":8080"); err!= nil{
		fmt.Println("Failed to start server", err)
	}
}