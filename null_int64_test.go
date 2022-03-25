package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/twharmon/gosql"
)

func TestNullInt64Valid(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `gosql:"primary"`
		Name gosql.NullInt64
	}
	control := T{
		ID:   5,
		Name: gosql.NullInt64{Valid: true, Int64: 5},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestNullInt64NotValid(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `gosql:"primary"`
		Name gosql.NullInt64
	}
	control := T{
		ID:   5,
		Name: gosql.NullInt64{Valid: false},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}
