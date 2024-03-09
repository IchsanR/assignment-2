package routers

import "github.com/gin-gonic/gin"

func StartServer() *gin.Engine {
	router := gin.Default()

	return router
}
