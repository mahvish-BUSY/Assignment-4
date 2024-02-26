package routes

import (
	controllers "assignment-4/Controllers"

	"github.com/gin-gonic/gin"
)

func BranchRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.GET("/branches/",controllers.GetAllBranches())
	incomingRoutes.GET("/branches/:branch_id",controllers.GetBranchDetails())
	incomingRoutes.POST("/branches",controllers.CreateBranch())
	incomingRoutes.PUT("/branches/:branch_id",controllers.UpdateBranch())
	incomingRoutes.DELETE("/branches/:branch_id",controllers.DeleteBranch())
}
