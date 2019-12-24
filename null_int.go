package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// NullInt .
type NullInt struct {
	Int   int
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullInt) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case int:
		n.Int = value.(int)
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 32)
		if err != nil {
			return err
		}
		n.Int = int(i)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int, nil
}

// MarshalJSON .
func (n NullInt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullInt) UnmarshalJSON(data []byte) error {
	var i *int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		n.Int, n.Valid = *i, true
	} else {
		n.Valid = false
	}
	return nil
}
