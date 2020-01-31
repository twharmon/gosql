package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSelectQueryOne(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, db.Register(T{}))
	control := T{
		ID:   5,
		Name: "foo",
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestSelectQueryMany(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, db.Register(T{}))
	control := []*T{
		&T{
			ID:   5,
			Name: "foo",
		},
		&T{
			ID:   6,
			Name: "bar",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Name)
	}
	mock.ExpectQuery(`^select \* from t limit 10$`).WillReturnRows(rows)
	var test []*T
	check(t, db.Select("*").Limit(10).Get(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control); i++ {
		equals(t, control[i], test[i])
	}
}

func TestSelectQueryManyValues(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, db.Register(T{}))
	control := []T{
		T{
			ID:   5,
			Name: "foo",
		},
		T{
			ID:   6,
			Name: "bar",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Name)
	}
	mock.ExpectQuery(`^select \* from t limit 10$`).WillReturnRows(rows)
	var test []T
	check(t, db.Select("*").Limit(10).Get(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control); i++ {
		equals(t, control[i], test[i])
	}
}
