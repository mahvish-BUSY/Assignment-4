package models

import (
	"time"
)

type Customer struct {
	CustId  uint `pg:",pk" json:"cust_id"`
	Name    string	`pg:",notnull" json:"name"`
	PAN     string	`pg:",notnull" json:"pan"`
	DOB     time.Time	`pg:",notnull" json:"dob"`
	Phone   int		`pg:",notnull" json:"phone"`
	Address string	`pg:",notnull" json:"address"`

	BranchID uint    `json:"branchId" pg:",on_delete:CASCADE"`
	Branch   *Branch `pg:"rel:has-one"`

	Account []*Account `pg:"many2many:customer_to_accounts"`
}
