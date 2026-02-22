package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/database"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/internal/ai"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/routes"
)

func init() {

	//Load environment variable
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
}

func main(){

	if os.Getenv("MONGODB_URI") == "" {
		fmt.Println("MONGO_URI not set, exiting...")
		os.Exit(1)
	}

	if os.Getenv("NVIDIA_API_KEY") == "" {
		fmt.Println("NVIDIA_API_KEY not set, exiting...")
		os.Exit(1)
	}

	//connect MongoDB 
	database.Connect()

	//Initialize controller dependencies
	controller.InitMovieController()
	controller.InitUserController()

	//LLM Model integration
	nvidiaClient := ai.NewNvidiaClient()
	aiService := ai.NewService(nvidiaClient)
	aiController := controller.InitAIController(aiService)

	//Create gin router
	router:= gin.Default()

	router.GET("/hello", func(c *gin.Context){
		c.String(200, "Hello, Magic Stream Movies!")
	})

	api := router.Group("/api/v1")
	routes.SetupUnProtectedRoutes(api)
	routes.SetupProctectedRoutes(api, aiController)

	//Load API key
	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		fmt.Println("NVIDIA_API_KEY not set")
		os.Exit(1)
	}

	if err:= router.Run(":8080"); err!= nil{
		fmt.Println("Failed to start server", err)
	}
}