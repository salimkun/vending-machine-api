package main

import (
	"fmt"

	_ "github.com/salimkun/vending-machine-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/salimkun/vending-machine-api/database"
	"github.com/salimkun/vending-machine-api/handler"
	"github.com/salimkun/vending-machine-api/middleware"
	"github.com/salimkun/vending-machine-api/payload"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router
// @title Vending Machine
// @version 0.0.1
// @description Vending Machine
// @schemes http
// @BasePath /
// @securityDefinitions.apikey  ApiJwtToken
// @in                          header
// @name                        Authorization
// Router list
func main() {
	fmt.Print("Code is ", " starting.\n")
	router := gin.Default()

	database.Connect()
	database.Migrate()

	api := router.Group("/api")
	{
		api.POST("/token", handler.GenerateToken)
		api.POST("/user/register", handler.RegisterUser)
		product := api.Group("/product").Use(middleware.AuthJwt())
		{
			product.GET("/", handler.GetAllProduct)
			product.GET("/:id", handler.GetDetailProductByID)
			product.PATCH("/:id", handler.UpdateProduct)
			product.DELETE("/:id", handler.DeleteProduct)
			product.POST("/", handler.CreateProduct)
		}

		buy := api.Group("/buy").Use(middleware.AuthJwt())
		{
			buy.POST("/product", handler.BuyProduct)
		}

	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404,
			payload.ResponseError{
				Success: false,
				Message: "error",
				Error:   "Page not found",
			})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run("localhost:8030")

}
