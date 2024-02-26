package models

import (
	"time"
)

type Transaction struct {
	ID       int `pg:"id,pk" json:"id"`
	Mode     string	`json:"mode"`
	RcvAccNo uint		`json:"rcv_acc_no"`
	TrDate   time.Time	`pg:"default:CURRENT_TIMESTAMP"`
	
	Amount   float64	`json:"amount"`

	AccId   uint	`json:"acc_id" pg:",on_delete:CASCADE"`
	Account *Account `pg:"rel:has-one"`
}
