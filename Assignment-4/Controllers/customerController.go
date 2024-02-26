package controllers

import (
	database "assignment-4/Database"
	models "assignment-4/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//gives customer details along with his accounts
func GetCustomerDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve the customer id from url
		custIdStr := c.Param("cust_id")
		custId,err := strconv.ParseUint(custIdStr,10,64)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"message":"Failed to retrieve the customer id",
			})
			return 
		}

		//retrieve the customer details
		//var customer models.Customer
		db := database.ReturnDBInstance()

		customer,selErr:= GetCustDetails(db,uint(custId))
		if selErr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":"Failed to retrieve the customer details",
			})
			return 
		}

		//return response
		c.JSON(http.StatusOK, gin.H{
			"message":"Customer details retrieved",
			"details":customer,
		})
		
	}
}

func UpdateCustomeDetails() gin.HandlerFunc {
	return func( c *gin.Context){

		//retrieve the details from request payload
		var cust models.Customer
		err:=c.ShouldBindJSON(&cust)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message":"Failed to retrieve details to update",
			})
			return
		}

		//update the details
		db := database.ReturnDBInstance()

		//begin transaction
		tx,txErr := db.Begin()
		if txErr != nil{
			tx.Close()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":"Failed to start transaction",
			})
			return
		}
		
		//update details
		upErr := updateCustDetails(tx,cust)
		if upErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":"Failed to update details",
				
			})
			return
		}
		tx.Commit()

		//return appropriate response
		c.JSON(http.StatusOK, gin.H{
			"message":"Details updated successfully",
		})
		
	}
}