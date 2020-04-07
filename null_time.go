package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullTime holds an time.Time value that might be null in the
// database.
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	n.Time = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

// MarshalJSON implements the Marshaler interface.
func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the Unmarshaler interface.
func (n *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		n.Time, n.Valid = *t, true
	} else {
		n.Valid = false
	}
	return nil
}
