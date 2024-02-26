package routes

import (
	controllers "assignment-4/Controllers"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(incomingRoutes *gin.Engine){
	
	incomingRoutes.GET("/transaction/:id", controllers.ViewTransactionDetails())
	incomingRoutes.POST("/transaction/transfer", controllers.TransferFunds())
	incomingRoutes.GET("/transaction/search",controllers.SearchTransaction())

}