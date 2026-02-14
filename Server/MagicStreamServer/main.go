package main

import (
	"fmt"
	"log"

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

	if err:= router.Run(":8080"); err!= nil{
		fmt.Println("Failed to start server", err)
	}
}