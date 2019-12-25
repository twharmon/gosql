package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// NullFloat32 .
type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

// Scan implements the Scanner interface.
func (n *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case float32:
		n.Float32 = value.(float32)
	case []byte:
		f, err := strconv.ParseFloat(string(value.([]byte)), 32)
		if err != nil {
			return err
		}
		n.Float32 = float32(f)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullFloat32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return float64(n.Float32), nil
}

// MarshalJSON .
func (n NullFloat32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullFloat32) UnmarshalJSON(data []byte) error {
	var i *float32
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		n.Float32, n.Valid = *i, true
	} else {
		n.Valid = false
	}
	return nil
}
