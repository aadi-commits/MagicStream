package main

import (
	"fmt"

	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/gin-gonic/gin"
)

func main(){
	
	router:= gin.Default()

	router.GET("/hello", func(c *gin.Context){
		c.String(200, "Hello, Magic Stream Movies!")
	})

	router.GET("/movies", controller.GetMovies())

	if err:= router.Run(":8080"); err!= nil{
		fmt.Println("Failed to start server", err)
	}
}