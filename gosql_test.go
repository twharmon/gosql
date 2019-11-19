package gosql_test

import (
	"testing"

	"github.com/twharmon/gosql"
)

func TestRegister(t *testing.T) {
	assertPanic(t, "should panic if non struct passed to Register", func() {
		m := make(map[string]string)
		gosql.Register(m)
	})

	assertPanic(t, "should panic if struct registered twice", func() {
		type T struct{}
		gosql.Register(T{}, T{})
	})
}
