package models

import "github.com/google/uuid"

type Branch struct {
	ID uint	`json:"id"`
	Address  string
	IFSCCode uuid.UUID	`pg:"type:uuid,default:uuid_generate_v4()"`

	BankID  uint	`pg:",on_delete:CASCADE"`
	Bank  *Bank      `pg:"rel:has-one"`

	Account []*Account `pg:"rel:has-many"`
	Customer []*Customer `pg:"rel:has-many"`
}
