package model

import (
	"time"

	"github.com/google/uuid"
)

type RoutinePayment struct {
	Id               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	RoutineExpenseId uuid.UUID
	Amount           float64
	DueDate          time.Time
	IsPaid           bool
	TenureNo         int64
	PaidDate         time.Time
	AuditColumn
}
