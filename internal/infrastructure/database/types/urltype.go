package types

import (
	"database/sql/driver"
	"fmt"
)

type URLType string

const (
	URLTypeYoutube    URLType = "youtube"
	URLTypeSoundCloud URLType = "soundcloud"
	URLTypeGithub     URLType = "github"
	URLTypeSketchfab  URLType = "sketchfab"
	URLTypeUnityroom  URLType = "unityroom"
	URLTypeOther      URLType = "other"
)

// Value はdatabase/sql/driver.Valuerインターフェースを実装
func (v URLType) Value() (driver.Value, error) {
	return string(v), nil
}

// Scan はsql.Scannerインターフェースを実装
func (ut *URLType) Scan(value interface{}) error {
	if value == nil {
		*ut = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		*ut = URLType(s)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into URLType", value)
	}
}
