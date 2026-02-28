package routes

import (
	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.RouterGroup){
	router.GET("/movies", controller.GetMovies())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
}