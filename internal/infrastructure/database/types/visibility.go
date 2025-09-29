package types

import (
	"database/sql/driver"
	"fmt"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
	VisibilityDraft   Visibility = "draft"
)

// Value はdatabase/sql/driver.Valuerインターフェースを実装
func (v Visibility) Value() (driver.Value, error) {
	return string(v), nil
}

// Scan はsql.Scannerインターフェースを実装
func (v *Visibility) Scan(value interface{}) error {
	if value == nil {
		*v = ""
		return nil
	}
	switch s := value.(type) {
	case []byte:
		*v = Visibility(string(s))
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Visibility", value)
	}
}
