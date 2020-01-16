package gosql_test

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/twharmon/gosql"
)

type fataler interface {
	Fatal(...interface{})
}

func equals(t *testing.T, a interface{}, b interface{}) {
	if reflect.TypeOf(a).Kind() == reflect.Ptr && reflect.TypeOf(b).Kind() == reflect.Ptr {
		a = reflect.ValueOf(a).Elem().Interface()
		b = reflect.ValueOf(b).Elem().Interface()
	}
	if a != b {
		t.Fatalf("expected %v to equal %v", a, b)
	}
}

func check(f fataler, err error) {
	if err != nil {
		f.Fatal(err)
	}
}

func getMockDB() (*gosql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return gosql.Conn(db), mock, err
}

func getSQLiteDB(f fataler, q string) *gosql.DB {
	os.Remove("/tmp/foo.db")
	sqliteDB, err := sql.Open("sqlite3", "/tmp/foo.db")
	check(f, err)
	sqliteDB.Exec(q)
	return gosql.Conn(sqliteDB)
}
