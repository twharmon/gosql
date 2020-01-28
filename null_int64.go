package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// NullInt64 .
type NullInt64 struct {
	Int64 int64
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case int64:
		n.Int64 = value.(int64)
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return err
		}
		n.Int64 = i
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

// MarshalJSON .
func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullInt64) UnmarshalJSON(data []byte) error {
	var i64 *int64
	if err := json.Unmarshal(data, &i64); err != nil {
		return err
	}
	if i64 != nil {
		n.Int64, n.Valid = *i64, true
	} else {
		n.Valid = false
	}
	return nil
}
