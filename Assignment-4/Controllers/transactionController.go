package controllers

import (
	database "assignment-4/Database"
	models "assignment-4/Models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ViewTransactionDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve the transaction id from request url
		transIdStr := c.Param("id")
		transId, err := strconv.ParseUint(transIdStr, 10, strconv.IntSize)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"message":"Failed to retrieve  transaction id from url",
			})
			return
		}

		//retrieve the transaction details
		db:= database.ReturnDBInstance() 

		if fetchedTrans, selErr := ViewTransDetails(db,uint(transId)); selErr == nil{
			//return appropriate response
			c.JSON(http.StatusOK, gin.H{
				"message":"Details retrieved successfully",
				"transaction details": fetchedTrans,
			})
		}else{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":"Failed to retrieve  transaction details",
			})
		}
			
	}
}

func TransferFunds() gin.HandlerFunc {eep
	return func(c *gin.Context) {

		//retrieve data from request payload
		transferDetails := models.Transaction{}
		if err := c.ShouldBindJSON(&transferDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		db := database.ReturnDBInstance()

		//get sender and reciever accounts
		sender, err := GetAccDetails(db,transferDetails.AccId)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve sender's account details",
				"error":   err.Error(),
			})
			return
		}

		var receiver models.Account
		if selErr := db.Model(&receiver).Where("acc_no=?", transferDetails.RcvAccNo).Select(); selErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve reciever's account details",
			})
			return
		}
		//begin transaction
		tx, txErr := db.Begin()

		if txErr != nil {
			tx.Close()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}
		
		//record this in transactions table
		if insertErr := InsertTransaction(tx,transferDetails); insertErr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to insert record in transactions table",
				"error":insertErr.Error(),
			})
			return
		}
		
		sender.Balance -= transferDetails.Amount
		receiver.Balance += transferDetails.Amount
		//update the record in accounts table for both sender and reciever
		if err := UpdateAcc(tx,sender); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender's account balance"})
			return
		}
		
		if err := UpdateAcc(tx,receiver); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receiver's account balance"})
			return
		}
		

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
		//return appropriate response
		c.JSON(http.StatusOK, gin.H{
			"message": "Transfer completed successfully",
		})

	}
}

func SearchTransaction() gin.HandlerFunc {
	return func (c *gin.Context){
		//retrieve the start_date from query param
		startDateStr := c.Query("start_date")

		startDate, err := time.Parse(time.RFC3339,startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message":"Failed to retrieve the date",
			})
			return 
		}
		//get db instance
		db := database.ReturnDBInstance()
		//retrieve the transactions
		if transactions, selErr := SearchTrans(db,startDate); selErr == nil{
			//return the response
			c.JSON(http.StatusOK, gin.H{
				"message":"Transactions retrieved",
				"Transactions":transactions,
			})
		} else {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message":"Failed to retrieve the transactions",
			})
			
		}
	}
}
