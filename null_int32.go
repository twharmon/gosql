package gosql

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"strconv"
)

// NullInt32 .
type NullInt32 struct {
	Int32 int32
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullInt32) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch value.(type) {
	case int32:
		n.Int32 = value.(int32)
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 32)
		if err != nil {
			return err
		}
		n.Int32 = int32(i)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n.Int32))
	return b, nil
}

// MarshalJSON .
func (n NullInt32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *NullInt32) UnmarshalJSON(data []byte) error {
	var i64 *int32
	if err := json.Unmarshal(data, &i64); err != nil {
		return err
	}
	if i64 != nil {
		n.Int32, n.Valid = *i64, true
	} else {
		n.Valid = false
	}
	return nil
}
