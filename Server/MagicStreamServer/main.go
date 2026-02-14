package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/database"
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

	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.AddMovie())

	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())

	if err:= router.Run(":8080"); err!= nil{
		fmt.Println("Failed to start server", err)
	}
}