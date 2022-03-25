package gosql_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/twharmon/gosql"
)

func TestNullTimeValid(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `gosql:"primary"`
		Name gosql.NullTime
	}
	control := T{
		ID:   5,
		Name: gosql.NullTime{Valid: true, Time: time.Now()},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestNullTimeNotValid(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `gosql:"primary"`
		Name gosql.NullTime
	}
	control := T{
		ID:   5,
		Name: gosql.NullTime{Valid: false},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}
