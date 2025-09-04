package types

import (
	"database/sql/driver"
	"fmt"
)

type AssetType string

const (
	AssetTypeZip   AssetType = "zip"
	AssetTypeImage AssetType = "image"
	AssetTypeVideo AssetType = "video"
	AssetTypeMusic AssetType = "music"
	AssetTypeModel AssetType = "model"
)

// Value はdatabase/sql/driver.Valuerインターフェースを実装
func (at AssetType) Value() (driver.Value, error) {
	return string(at), nil
}

// Scan はsql.Scannerインターフェースを実装
func (at *AssetType) Scan(value interface{}) error {
	if value == nil {
		*at = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		*at = AssetType(s)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into AssetType", value)
	}
}
