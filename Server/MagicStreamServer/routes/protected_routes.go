package routes

import (
	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupProctectedRoutes(router *gin.RouterGroup, ctl *controller.AIController){

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())

	protected.GET("/movie/:imdb_id", controller.GetMovie())
	protected.POST("/addmovie", controller.AddMovie())
	protected.POST("/ai/infer", ctl.Infer())
}