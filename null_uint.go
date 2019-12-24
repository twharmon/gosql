package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// NullUint .
type NullUint struct {
	Uint  uint
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullUint) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case uint:
		n.Uint = value.(uint)
	case []byte:
		i, err := strconv.ParseUint(string(value.([]byte)), 10, 32)
		if err != nil {
			return err
		}
		n.Uint = uint(i)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUint) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint, nil
}

// MarshalJSON .
func (n NullUint) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Uint)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullUint) UnmarshalJSON(data []byte) error {
	var i *uint
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		n.Uint, n.Valid = *i, true
	} else {
		n.Valid = false
	}
	return nil
}
