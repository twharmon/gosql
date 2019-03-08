package gosql_test

import (
	"testing"

	"github.com/twharmon/gosql"
)

func TestRegister(t *testing.T) {
	assertPanic(t, "should panic if non struct passed to Register", func() {
		gosql.Register(&User{})
	})

	assertPanic(t, "should panic if struct registered twice", func() {
		gosql.Register(User{}, User{})
	})

	assertPanic(t, "should panic if first field not ID", func() {
		type Test struct {
			Name string
			ID   int
		}
		gosql.Register(Test{})
	})

	assertPanic(t, "should panic if ID not of type int64", func() {
		type Test struct {
			ID   int
			Name string
		}
		gosql.Register(Test{})
	})
}
