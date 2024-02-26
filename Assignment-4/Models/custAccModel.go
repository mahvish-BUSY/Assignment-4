package models

type CustomerToAccount struct {
	ID uint

	AccId   uint     `pg:",on_delete:CASCADE"`
	Account *Account `pg:"rel:has-one"`

	CustId   uint      `pg:",on_delete:CASCADE"`
	Customer *Customer `pg:"rel:has-one"`
}
