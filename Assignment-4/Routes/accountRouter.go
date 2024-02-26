package routes

import (
	controllers "assignment-4/Controllers"

	"github.com/gin-gonic/gin"
)

func AccountRoutes( incomingRoutes *gin.Engine){

	incomingRoutes.POST("/accounts/open", controllers.OpenAccount())
	incomingRoutes.POST("/accounts/close/:acc_id",controllers.CloseAccount())
	incomingRoutes.GET("/accounts/:acc_id",controllers.GetAccountDetails())
	incomingRoutes.POST("/accounts/deposit",controllers.DepositFunds())
	incomingRoutes.POST("/accounts/withdraw",controllers.WithdrawFunds())
	incomingRoutes.GET("/accounts/:acc_id/transactionhistory",controllers.GetTransactionDetails())
	
}