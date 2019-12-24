package gosql

import (
	"database/sql/driver"
	"encoding/json"
)

// NullString .
type NullString struct {
	String string
	Valid  bool
}

// Scan implements the Scanner interface.
func (n *NullString) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	switch value.(type) {
	case string:
		n.String, n.Valid = value.(string), true
	case []byte:
		n.String, n.Valid = string(value.([]byte)), true
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullString) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

// MarshalJSON .
func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		n.String, n.Valid = *s, true
	} else {
		n.Valid = false
	}
	return nil
}
