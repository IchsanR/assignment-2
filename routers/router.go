package routers

import (
	"assigntment2/controllers"
	"assigntment2/repository"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func StartServer(db *sql.DB) *gin.Engine {
	router := gin.Default()

	orderRepo := repository.NewRepo(db)
	orderController := controllers.NewOrderController(orderRepo)

	router.POST("/orders", orderController.CreateOrder)

	return router
}
