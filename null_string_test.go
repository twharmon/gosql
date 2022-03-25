package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/twharmon/gosql"
)

func TestNullStringValid(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name gosql.NullString
	}
	control := T{
		ID:   5,
		Name: gosql.NullString{Valid: true, String: "foo"},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestNullStringNotValid(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name gosql.NullString
	}
	control := T{
		ID:   5,
		Name: gosql.NullString{Valid: false},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}
