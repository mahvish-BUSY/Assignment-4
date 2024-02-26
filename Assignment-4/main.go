package main

import (
	database "assignment-4/Database"
	models "assignment-4/Models"
	routes "assignment-4/Routes"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var Pg_db *pg.DB

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*models.CustomerToAccount)(nil))
}
func main() {

	Pg_db = database.DBConnection()
	createErr := database.CreateTables(Pg_db)
	if createErr != nil {
		panic(createErr)
	}

	router := gin.Default()

	routes.AccountRoutes(router)
	routes.BankRoutes(router)
	routes.BranchRoutes(router)
	routes.CustomerRoutes(router)
	routes.TransactionRoutes(router)
	router.Run()
}
