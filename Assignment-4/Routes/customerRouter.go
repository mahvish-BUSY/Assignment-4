package routes

import (
	controllers "assignment-4/Controllers"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.GET("/customers/:cust_id", controllers.GetCustomerDetails())
	incomingRoutes.PUT("/customers/update",controllers.UpdateCustomeDetails())
}