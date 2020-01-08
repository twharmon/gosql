package gosql

import (
	"testing"
)

func TestRegister(t *testing.T) {
	m := make(map[string]string)
	assertErr(t, "should return error if non struct passed to Register", Register(m))

	type T struct {
		ID uint64 `gosql:"primary"`
	}
	assertErr(t, "should return error if struct registered twice", Register(T{}, T{}))
}
