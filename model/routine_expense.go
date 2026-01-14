package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RoutineExpense struct {
	Id               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	InitialAmount    *float64
	TotalAmount      *float64
	Tenure           int64
	StartDate        time.Time
	EndDate          time.Time
	IsActive         bool
	RemainingPayment int64
	RemainingAmount  float64
	Object           int64
	Notes            string
	PaidOffDate      *time.Time
	Type             RoutineExpenseTypeEnum
	AuditColumn
	RoutinePayments []RoutinePayment
}

type RoutineExpenseTypeEnum string

const (
	SUBSCRIPTION RoutineExpenseTypeEnum = "subscription"
	INSTALLMENT  RoutineExpenseTypeEnum = "installment"
	OTHER        RoutineExpenseTypeEnum = "other"
)

func (p *RoutineExpenseTypeEnum) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*p = RoutineExpenseTypeEnum(string(v))
	case string:
		*p = RoutineExpenseTypeEnum(v)
	default:
		return fmt.Errorf("unsupported Scan type for RoutineExpenseTypeEnum: %T", value)
	}

	return nil
}

func (p RoutineExpenseTypeEnum) Value() (driver.Value, error) {
	return string(p), nil
}
