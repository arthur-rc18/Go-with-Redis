package routes

import (
	"net/http"

	"github.com/arthur-rc18/Go-Redis/handlers"
	"github.com/gin-gonic/gin"
)

func StartRoutes(routes *gin.Engine) {

	routes.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "Up and running..."})
	})

	blockRoutes := routes.Group("/blocks")
	{
		blockRoutes.GET("", handlers.GetBlocksData)
		blockRoutes.GET("/:id", handlers.GetBlockByID)
		blockRoutes.GET("/tree/:id", handlers.GetTreeID)
		blockRoutes.DELETE("/:id", handlers.DeleteBlockByID)
		blockRoutes.PUT("/:id", handlers.UpdateBlockByID)
		blockRoutes.POST("/", handlers.CreateBlock)

	}

	routes.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})

}
