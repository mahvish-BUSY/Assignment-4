package models

type Account struct {
	AccID   uint `pg:",pk" json:"acc_id"` //uuid.UUID `pg:"type:uuid,pk"`
	AccNo   uint `pg:",unique,notnull,default:nextval('account_number_sequence')"`
	Balance float64
	AccType string

	BranchID uint    `pg:",on_delete:CASCADE"`
	Branch   *Branch `pg:"rel:has-one"`

	Transaction []*Transaction `pg:"rel:has-many"`

	Customers []*Customer `pg:"many2many:customer_to_accounts"`
}
