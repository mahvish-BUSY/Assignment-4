package controllers

import (
	database "assignment-4/Database"
	models "assignment-4/Models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WithdrawOrDeposit struct {
	AccID  uint    `json:"acc_id"`
	Amount float64 `json:"amount"`
	Mode string `json:"mode"`
}

type OpenCustAcc struct {
	Customers []struct {
		Name    string    `json:"name"`
		PAN     string    `json:"pan"`
		DOB     time.Time `json:"dob"`
		Phone   int       `json:"phone"`
		Address string    `json:"address"`
	} `json:"customers"`
	BranchID uint    `json:"branch_id"`
	Balance  float64 `json:"balance"`
	AccType  string  `json:"acc_type"`
}

func OpenAccount() gin.HandlerFunc {
	return func(c *gin.Context) {

		//retrieve data from request paylod
		var inputDetails OpenCustAcc
		err := c.ShouldBindJSON(&inputDetails)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Failed to parse request payload",
			})
			return
		}
		//get db instance
		db := database.ReturnDBInstance()

		//begin a transaction
		tx, err := db.Begin()
		if err != nil {
			tx.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		var custIds []uint
		
		//insert record in customer table
		// custIds []uint will be used for mapping
		custIds, custErr := SaveCustomers(tx, inputDetails)
		if custErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": custErr.Error(),
			})
			return
		}

		//an instance of Account struct
		acc := &models.Account{
			Balance:  inputDetails.Balance,
			AccType:  inputDetails.AccType,
			BranchID: inputDetails.BranchID,
		}
		//insert record in accounts table
		accId, accErr := SaveAccount(tx, acc)
		if accErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": accErr.Error(),
			})
			return
		}
		
		//establish relation between accounts and customers table
		if mapErr := SaveCustAcc(tx, custIds, accId); mapErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   mapErr.Error(),
				"message": "Failed to insert record in accounts table",
			})
			return
		}

		//return the response
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
		//return success response
		c.JSON(http.StatusCreated, gin.H{
			"message": "Account opened successfully",
		})
	}
}

func CloseAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve acc_id
		accIdStr := c.Param("acc_id")
		accID, err := strconv.ParseUint(accIdStr, 10, strconv.IntSize)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		db := database.ReturnDBInstance()

		tx, txErr := db.Begin()
		if txErr != nil {
			
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}
		
		//retrieve corresponding customer details
		// accCust is a []models.CustomerToAccount
		accCust, selErr := GetCustAcc(tx,uint(accID))
		
		if selErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve account data from mapping table",
			})
			return
		}

		//iterating over this slice to delete corresponding customer details
		if delErr := DeleteCust(tx,accCust); delErr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":delErr.Error(),
			})
			return 
		}
		
		if delErr := DeleteAcc(tx,uint(accID)); delErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to delete account details from account table",
				"error":   delErr.Error(),
			})
			return
		}

		//commit the transaction
		tx.Commit()
		//return response of success
		c.JSON(http.StatusOK, gin.H{
			"message": "Account closed successfully",
		})

	}
}

func GetAccountDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve acc_id from url
		accIdStr := c.Param("acc_id")
		accId, err := strconv.ParseUint(accIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Retrieve the account details
		db := database.ReturnDBInstance()

		if account, selectErr := GetAccDetails(db, uint(accId)); selectErr == nil {

			// Return the account details
			c.JSON(http.StatusOK, gin.H{
				"message": "Account details retrieved successfully",
				"account": account,
			})
		} else {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve account details",
			})
		}
	}
}

func WithdrawFunds() gin.HandlerFunc {
	return func(c *gin.Context) {
		//to store request payload
		var withdrawData WithdrawOrDeposit

		if err := c.ShouldBindJSON(&withdrawData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//get db instance
		db := database.ReturnDBInstance()

		//retrieve the account record from accounts table
		account, err := GetAccDetails(db,withdrawData.AccID)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}
		//open a transaction
		tx, txErr := db.Begin()

		if txErr != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}
		
		account.Balance -= withdrawData.Amount

		//record this in transaction table
		// Create a new transaction record for deposit
		transaction := models.Transaction{
			Mode:     withdrawData.Mode,
			RcvAccNo: account.AccNo,
			Amount:   withdrawData.Amount,
			AccId:    account.AccID,
		}
		//add record in transaction table
		if insertErr := InsertTransaction(tx,transaction); insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
			return
		}
		
		// Update the account record with the new balance
		if err := UpdateAcc(tx,account); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
			return
		}
		// commit the transaction
		tx.Commit()

		//return appropriate response
		c.JSON(http.StatusOK, gin.H{
			"message":     "Amount withdrawn successfully",
			"new Balance": account.Balance,
		})

	}
}

func DepositFunds() gin.HandlerFunc {
	return func(c *gin.Context) {

		//to store request payload
		var depositData WithdrawOrDeposit

		if err := c.ShouldBindJSON(&depositData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//get db instance
		db := database.ReturnDBInstance()

		//retrieve the account record from accounts table
		account, err := GetAccDetails(db,depositData.AccID)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}

		//open a transaction
		tx, txErr := db.Begin()

		if txErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}
		
		// Add the deposit amount to the current balance
		account.Balance += depositData.Amount

		//record this in transaction table
		// Create a new transaction record for the deposit
		transaction := models.Transaction{
			Mode:     depositData.Mode,
			RcvAccNo: account.AccNo,
			Amount:   depositData.Amount,
			AccId:    account.AccID,
		}
		//add record in transaction table
		if insertErr := InsertTransaction(tx,transaction); insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
			return
		}

		// Update the account record with the new balance
		if err := UpdateAcc(tx,account); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
			return
		}

		// commit the transaction
		tx.Commit()

		//return appropriate response
		c.JSON(http.StatusOK, gin.H{
			"message":     "Amount deposited successfully",
			"new Balance": account.Balance,
		})

	}
}

func GetTransactionDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve acc_id from request url
		accIdStr := c.Param("acc_id")
		accId, err := strconv.ParseUint(accIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to retrieve account id from request url",
			})
			return
		}

		//retrieve all the transactions related to that account
		db := database.ReturnDBInstance()

		if account, selErr := GetTransDetailsOfAcc(db, uint(accId)); selErr == nil {
			//return the appropriate response
			c.JSON(http.StatusOK, gin.H{
				"message":             "Transaction details retrieved successfully",
				"transaction details": account,
			})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve transaction details related to thi account",
			})

		}

	}
}
