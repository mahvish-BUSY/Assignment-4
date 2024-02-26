package controllers

import (
	database "assignment-4/Database"
	models "assignment-4/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllBanks() gin.HandlerFunc {

	return func(c *gin.Context) {
		//retrieve all records from bank table
		var banks []*models.Bank
		db := database.ReturnDBInstance()

		banks, selErr := GetBanks(db)

		if selErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error while retrieving the details of banks",
			})
			return
		}
		//return the response
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Details retrieved",
			"banks":   banks,
		})
	}

}

// returns a specific bank detail along with its branches
func GetBankDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve the bankId
		bankIdStr := c.Param("bank_id")
		bankId, err := strconv.ParseUint(bankIdStr, 10, strconv.IntSize)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//retrieve the bank details along with its branches
		var bank models.Bank
		db := database.ReturnDBInstance()

		bank, selectErr := GetBankById(db, uint(bankId))

		if selectErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": selectErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Bank details retrieved",
			"details": bank,
		})

	}
}

func UpdateBankDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get id from url parameters
		bankIdStr := c.Param("bank_id")
		bankID, err := strconv.ParseUint(bankIdStr, 10, strconv.IntSize)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//retrieve the data from json which is to be updated
		bank := &models.Bank{}

		if err := c.ShouldBindJSON(bank); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//get db instance
		db := database.ReturnDBInstance()
		//update the record
		updateBankId, updateErr := UpdateBank(db, bank, uint(bankID))
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while updating the record",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Bank Id": updateBankId,
			"message": "Record updated successfully",
		})

	}
}

func CreateBank() gin.HandlerFunc {

	return func(c *gin.Context) {
		//an empty instance of bank model
		bank := &models.Bank{}

		if err := c.ShouldBindJSON(bank); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//insert the record in database
		
		db := database.ReturnDBInstance()
		newBankId , err := CreateNewBank(db,bank)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some db error in creating bank",
			})
			return
		}
		//return response
		c.JSON(http.StatusOK, gin.H{
			"message":         "Bank created successfully",
			"inserted Bank Id": newBankId,
		})
	}

}

func DeleteBank() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve the bankId
		bankIdStr := c.Param("bank_id")
		bankId, err := strconv.ParseUint(bankIdStr, 10, strconv.IntSize)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		db := database.ReturnDBInstance()
		//deleting a bank will be a transaction

		tx, txErr := db.Begin()
		if txErr != nil {
			tx.Close()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": txErr.Error(),
			})
			return
		}

		// Delete the bank record
		if deleteErr := DeleteExistingBank(tx,uint(bankId));deleteErr == nil {

			//all good then commit and return response
			if commitErr := tx.Commit(); commitErr == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "Bank deleted successfully",
				})
			}
		}else{

			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{
				"message": "Failed to commit transaction",
			})
			return
		}

	}
}
