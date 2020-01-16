package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/twharmon/gosql"
)

func TestSelectQueryOne(t *testing.T) {
	type SelectQueryOneModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(SelectQueryOneModel{}))
	control := SelectQueryOneModel{
		ID:   5,
		Name: "foo",
	}
	db, mock, err := getMockDB()
	check(t, err)
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from select_query_one_model where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test SelectQueryOneModel
	check(t, db.Select("*").Where("id = ?", control.ID).To(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestSelectQueryMany(t *testing.T) {
	type SelectQueryManyModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(SelectQueryManyModel{}))
	control := []*SelectQueryManyModel{
		&SelectQueryManyModel{
			ID:   5,
			Name: "foo",
		},
		&SelectQueryManyModel{
			ID:   6,
			Name: "bar",
		},
	}
	db, mock, err := getMockDB()
	check(t, err)
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Name)
	}
	mock.ExpectQuery(`^select \* from select_query_many_model limit 10$`).WillReturnRows(rows)
	var test []*SelectQueryManyModel
	check(t, db.Select("*").Limit(10).To(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control); i++ {
		equals(t, control[i], test[i])
	}
}

func TestSelectQueryManyValues(t *testing.T) {
	type SelectQueryManyValuesModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(SelectQueryManyValuesModel{}))
	control := []SelectQueryManyValuesModel{
		SelectQueryManyValuesModel{
			ID:   5,
			Name: "foo",
		},
		SelectQueryManyValuesModel{
			ID:   6,
			Name: "bar",
		},
	}
	db, mock, err := getMockDB()
	check(t, err)
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Name)
	}
	mock.ExpectQuery(`^select \* from select_query_many_values_model limit 10$`).WillReturnRows(rows)
	var test []SelectQueryManyValuesModel
	check(t, db.Select("*").Limit(10).To(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control); i++ {
		equals(t, control[i], test[i])
	}
}
