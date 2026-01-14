package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type Account struct {
	Id            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	OwnerName     string
	AccountNumber string
	AccountType   AccountTypeEnum `gorm:"type:account_type_enum"`
	Balance       float64
	ProviderName  *string
	Purpose       *string
	Transactions  []Transaction
	AuditColumn
}

type AccountTypeEnum string

const (
	BANK_ACCOUNT AccountTypeEnum = "bank_account"
	EWALLET      AccountTypeEnum = "ewallet"
	CREDIT_CARD  AccountTypeEnum = "credit_card"
	CASH         AccountTypeEnum = "cash"
	RDN          AccountTypeEnum = "rdn"
)

func (p *AccountTypeEnum) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*p = AccountTypeEnum(string(v))
	case string:
		*p = AccountTypeEnum(v)
	default:
		return fmt.Errorf("unsupported Scan type for AccountTypeEnum: %T", value)
	}

	return nil
}

func (p AccountTypeEnum) Value() (driver.Value, error) {
	return string(p), nil
}
