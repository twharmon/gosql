package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// NullFloat64 .
type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

// Scan implements the Scanner interface.
func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case float64:
		n.Float64 = value.(float64)
	case []byte:
		f, err := strconv.ParseFloat(string(value.([]byte)), 64)
		if err != nil {
			return err
		}
		n.Float64 = f
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

// MarshalJSON .
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	var i *float64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		n.Float64, n.Valid = *i, true
	} else {
		n.Valid = false
	}
	return nil
}
