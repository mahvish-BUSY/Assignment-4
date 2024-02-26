package routes

import (
	controllers "assignment-4/Controllers"

	"github.com/gin-gonic/gin"
)

func BankRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/bank",controllers.GetAllBanks())
	incomingRoutes.GET("/bank/:bank_id", controllers.GetBankDetails())
	incomingRoutes.PUT("/bank/:bank_id",controllers.UpdateBankDetails())
	incomingRoutes.POST("/bank",controllers.CreateBank())
	incomingRoutes.DELETE("/bank/:bank_id",controllers.DeleteBank())
}
