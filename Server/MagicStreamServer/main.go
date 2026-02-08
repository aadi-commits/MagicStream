package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/database"
)

func main(){
	
	//Load environment variable
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	//connect MongoDB 
	database.Connect()

	//Create gin router
	router:= gin.Default()

	router.GET("/hello", func(c *gin.Context){
		c.String(200, "Hello, Magic Stream Movies!")
	})

	router.GET("/movies", controller.GetMovies())

	if err:= router.Run(":8080"); err!= nil{
		fmt.Println("Failed to start server", err)
	}
}