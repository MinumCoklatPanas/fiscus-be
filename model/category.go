package model

import "github.com/google/uuid"

type Category struct {
	Id           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name         string
	Transactions []Transaction
	AuditColumn
}
