package controllers

import (
	database "assignment-4/Database"
	models "assignment-4/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllBranches() gin.HandlerFunc {
	return func(c *gin.Context) {

		var branches []*models.Branch
		db := database.ReturnDBInstance()

		branches ,selectErr := GetBranches(db)

		if selectErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": selectErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":"Branch details retrieved",
			"branch details":branches,
		})	
	}
}

func GetBranchDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

		branchIdStr := c.Param("branch_id")
		if branchIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please give the branch id",
			})
			return
		}

		branchId, err := strconv.ParseUint(branchIdStr, 10, strconv.IntSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		db := database.ReturnDBInstance()

		if branch, selectErr := GetBranchById(db,uint(branchId)); selectErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": selectErr.Error(),
			})

		}else{
			c.JSON(http.StatusOK, gin.H{
				"message":"Branch details retrieved",
				"branch details":branch,
			})
		}
		
	}
}

func CreateBranch() gin.HandlerFunc {
	return func(c *gin.Context) {

		var branch models.Branch

		//bind var branch to Json data
		if err := c.ShouldBindJSON(&branch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//check if the bank exists in banks table
		db := database.ReturnDBInstance()
		count, err := db.Model((*models.Bank)(nil)).Where("id=?", branch.BankID).Count()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error in checking bank existence",
			})
			return
		}
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Specified bank does not exist",
			})
			return
		}

		//begin a transaction and create a branch in branch table.
		tx, txErr := db.Begin()

		// Make sure to close transaction if something goes wrong.
		if txErr != nil {
			tx.Close()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error in establishing transaction",
				"error":   txErr.Error(),
			})
			return
		}

		defer tx.Rollback() // Rollback transaction if any error occurs

		//create new record in branch table.
		_, createErr := tx.Model(&branch).Insert()

		if createErr != nil {
			//rollback changes
			tx.Rollback()

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while creating the branch record",
				"error":   createErr.Error(),
			})

			return
		}

		//all good then commit and return response
		tx.Commit()

		c.JSON(http.StatusOK, gin.H{
			"message":         "Branch created successfully",
			"inserted record": branch,
		})

	}
}

func UpdateBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get id from url parameters
		branchIdStr := c.Param("branch_id")
		branchID, err := strconv.ParseUint(branchIdStr, 10, strconv.IntSize)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//retrieve the data from json which is to be updated
		branch := &models.Branch{}

		if err := c.ShouldBindJSON(branch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//get db instance
		db := database.ReturnDBInstance()

		//begin transaction
		tx,txErr := db.Begin()
		if txErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":"Failed to start transaction",
			})
		}

		updatedBranchId, err := UpdateBranchRec(tx,uint(branchID),branch)

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		//commit the transaction
		tx.Commit()
		c.JSON(http.StatusCreated, gin.H{
			"message":"Branch has been updated successfully",
			"Branch ID":updatedBranchId,
		})


	}
}

func DeleteBranch() gin.HandlerFunc {
	return func(c *gin.Context){
		//retrieve the branchId
		branchIdStr := c.Param("branch_id")
		branchId, parseErr := strconv.ParseUint(branchIdStr, 10, strconv.IntSize)

		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
			"error": parseErr.Error(),
			})
		return
		}

		db := database.ReturnDBInstance()
		//deleting a branch will be a transaction

		tx, txErr := db.Begin()
		if txErr != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"error": txErr.Error(),
			})
			return 
		}

		// Delete the bank record
		deleteErr := DeleteBranchById(tx,uint(branchId));

		if deleteErr != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":deleteErr.Error()})
			return
		}

		//all good then commit and return response
		tx.Commit()

		c.JSON(http.StatusOK, gin.H{
			"message": "Branch deleted successfully",
		})
	}
}
