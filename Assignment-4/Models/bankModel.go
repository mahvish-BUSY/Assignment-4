package models

type Bank struct {
	ID     uint      `json:"id"`//uuid.UUID `pg:"type:uuid,pk"`
	Name   string    `json:"name"`
	Branch []*Branch `pg:"rel:has-many"`
}
