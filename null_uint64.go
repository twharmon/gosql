package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"strconv"
)

// NullUint64 .
type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

// Scan implements the Scanner interface.
func (n *NullUint64) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case uint64:
		log.Println("uint64")
		n.Uint64 = value.(uint64)
	case []byte:
		log.Println("[]byte")
		i, err := strconv.ParseUint(string(value.([]byte)), 10, 64)
		if err != nil {
			return err
		}
		n.Uint64 = i
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint64, nil
}

// MarshalJSON .
func (n NullUint64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Uint64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullUint64) UnmarshalJSON(data []byte) error {
	var u64 *uint64
	if err := json.Unmarshal(data, &u64); err != nil {
		return err
	}
	if u64 != nil {
		n.Uint64, n.Valid = *u64, true
	} else {
		n.Valid = false
	}
	return nil
}
