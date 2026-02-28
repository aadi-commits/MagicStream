package routes

import (
	controller "github.com/aadi-commits/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupProctectedRoutes(router *gin.RouterGroup, ctl *controller.AIController){

	auth := router.Group("/")
	auth.Use(middlewares.AuthMiddleware())

	auth.GET("/movie/:imdb_id",
		middlewares.Authorize(1),
		controller.GetMovie(),
	)

	admin := auth.Group("/")
	admin.Use(middlewares.Authorize(2))

	admin.POST("/addmovie", controller.AddMovie())
	admin.POST("/ai/infer", ctl.Infer())
}