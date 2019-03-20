package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// NullBool .
type NullBool struct {
	Bool  bool
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case bool:
		n.Bool = value.(bool)
	case []byte:
		b, err := strconv.ParseBool(string(value.([]byte)))
		if err != nil {
			return err
		}
		n.Bool = b
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

// MarshalJSON .
func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		n.Bool, n.Valid = *b, true
	} else {
		n.Valid = false
	}
	return nil
}
