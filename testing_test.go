package gosql_test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/twharmon/gosql"
)

func equals(t *testing.T, a interface{}, b interface{}) {
	if reflect.TypeOf(a).Kind() == reflect.Ptr && reflect.TypeOf(b).Kind() == reflect.Ptr {
		a = reflect.ValueOf(a).Elem().Interface()
		b = reflect.ValueOf(b).Elem().Interface()
	}
	if a != b {
		t.Fatalf("expected %v to equal %v", a, b)
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func getMockDB() (*gosql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return gosql.Conn(db), mock, err
}
