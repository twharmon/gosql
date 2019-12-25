package gosql

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"strconv"
)

// NullUint32 .
type NullUint32 struct {
	Uint32 uint32
	Valid  bool
}

// Scan implements the Scanner interface.
func (n *NullUint32) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case uint32:
		n.Uint32 = value.(uint32)
	case []byte:
		i, err := strconv.ParseUint(string(value.([]byte)), 10, 32)
		if err != nil {
			return err
		}
		n.Uint32 = uint32(i)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUint32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, n.Uint32)
	return b, nil
}

// MarshalJSON .
func (n NullUint32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Uint32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullUint32) UnmarshalJSON(data []byte) error {
	var u64 *uint32
	if err := json.Unmarshal(data, &u64); err != nil {
		return err
	}
	if u64 != nil {
		n.Uint32, n.Valid = *u64, true
	} else {
		n.Valid = false
	}
	return nil
}
