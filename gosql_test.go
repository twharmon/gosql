package gosql_test

import (
	"testing"

	"github.com/twharmon/gosql"
)

func TestRegister(t *testing.T) {
	m := make(map[string]string)
	assertErr(t, "should return error if non struct passed to Register", gosql.Register(m))

	type T struct {
		ID uint64 `gosql:"primary"`
	}
	assertErr(t, "should return error if struct registered twice", gosql.Register(T{}, T{}))
}
