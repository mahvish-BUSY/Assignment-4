package database

import (
	models "assignment-4/Models"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var db *pg.DB
func DBConnection() *pg.DB {
	opts := pg.Options{
		Addr:     ":5432",
		User:     "app_user",
		Password: "app_password",
		Database: "app_database",
	}

	db = pg.Connect(&opts)

	if db == nil {
		log.Printf("Failed to Connect to database")
		os.Exit(100)
	}
	log.Printf("Connection to database successful")
	return db
}

func CreateTables(db *pg.DB) error {
	models := []interface{}{
		(*models.Bank)(nil),
		(*models.Branch)(nil),
		(*models.Account)(nil),
		(*models.Customer)(nil),
		(*models.CustomerToAccount)(nil),
		(*models.Transaction)(nil),

	}

	for _, model := range models {
		createTableErr := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})

		if createTableErr != nil {
			log.Printf("Error in creating table ,Reason :%v", createTableErr)
			return createTableErr
		}
	}
	return nil
}

func ReturnDBInstance() *pg.DB{
	return db
}
