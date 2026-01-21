package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id         int64 `gorm:"primaryKey;type:serial"`
	AccountId  uuid.UUID
	Amount     float64
	TrxDate    time.Time
	CategoryId uuid.UUID `gorm:"type:uuid"`
	Notes      *string
	TrxType    TrxTypeEnum `gorm:"type:trx_type_enum"`
	AuditColumn
}

type TrxTypeEnum string

const (
	INCOME  TrxTypeEnum = "income"
	EXPENSE TrxTypeEnum = "expense"
)

func (p *TrxTypeEnum) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*p = TrxTypeEnum(string(v))
	case string:
		*p = TrxTypeEnum(v)
	default:
		return fmt.Errorf("unsupported Scan type for TrxTypeEnum: %T", value)
	}

	return nil
}

func (p TrxTypeEnum) Value() (driver.Value, error) {
	return string(p), nil
}
