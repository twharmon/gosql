package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSelectQueryOne(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
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
		ID   int `idx:"primary"`
		Name string
	}
	control := []*T{
		{
			ID:   5,
			Name: "foo",
		},
		{
			ID:   6,
			Name: "bar",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Name)
	}
	mock.ExpectQuery(`^select \* from t$`).WillReturnRows(rows)
	var test []*T
	check(t, db.Select("*").Get(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control); i++ {
		equals(t, control[i], test[i])
	}
}

func TestSelectQueryManyValues(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := []T{
		{
			ID:   5,
			Name: "foo",
		},
		{
			ID:   6,
			Name: "bar",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Name)
	}
	mock.ExpectQuery(`^select \* from t$`).WillReturnRows(rows)
	var test []T
	check(t, db.Select("*").Get(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control); i++ {
		equals(t, control[i], test[i])
	}
}

func TestSelectQueryManyLimit(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := []*T{
		{
			ID:   5,
			Name: "foo",
		},
		{
			ID:   6,
			Name: "bar",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Name)
	}
	mock.ExpectQuery(`^select \* from t$`).WillReturnRows(rows)
	var test []*T
	check(t, db.Select("*").Get(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control); i++ {
		equals(t, control[i], test[i])
	}
}

func TestSelectQueryManyValuesLimit(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := []T{
		{
			ID:   5,
			Name: "foo",
		},
		{
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

func TestSelectQueryOrWhere(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := T{
		ID:   5,
		Name: "foo",
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? or name = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).OrWhere("name = ?").Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestSelectQueryOrderBy(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := T{
		ID:   5,
		Name: "foo",
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t order by name asc limit 1$`).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").OrderBy("name asc").Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestSelectQueryOffset(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := []T{
		{
			ID:   6,
			Name: "bar",
		},
		{
			ID:   7,
			Name: "baz",
		},
		{
			ID:   5,
			Name: "foo",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for i, c := range control {
		if i != 2 {
			rows.AddRow(c.ID, c.Name)
		}
	}
	mock.ExpectQuery(`^select \* from t order by name asc limit 2 offset 1$`).WillReturnRows(rows)
	var test []T
	check(t, db.Select("*").OrderBy("name asc").Offset(1).Limit(2).Get(&test))
	check(t, mock.ExpectationsWereMet())
	for i := 0; i < len(control)-1; i++ {
		equals(t, control[i], test[i])
	}
}

func TestSelectQueryJoin(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := T{
		ID:   5,
		Name: "foo",
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select id, name from t join a on a\.t_id = t\.id order by name asc limit 1$`).WillReturnRows(rows)
	var test T
	check(t, db.Select("id", "name").Join("a on a.t_id = t.id").OrderBy("name asc").Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestSelectQueryErrNilPtr(t *testing.T) {
	db, _, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	var test *T
	if err := db.Select("*").Get(test); err == nil {
		t.Fatalf("expected err to be non nil")
	} else {
		contains(t, err.Error(), "nil")
	}
}

func TestSelectQueryErrNotPtr(t *testing.T) {
	db, _, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	var test T
	if err := db.Select("*").Get(test); err == nil {
		t.Fatalf("expected err to be non nil")
	} else {
		contains(t, err.Error(), "pointer")
	}
}

func TestSelectQueryErrNotStructOrSlice(t *testing.T) {
	db, _, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	var test map[string]string
	if err := db.Select("*").Get(&test); err == nil {
		t.Fatalf("expected err to be non nil")
	} else {
		contains(t, err.Error(), "struct")
	}
}

func TestSelectQueryCustomColumn(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int    `idx:"primary"`
		Name string `col:"first_name"`
	}
	control := T{
		ID:   5,
		Name: "foo",
	}
	rows := sqlmock.NewRows([]string{"id", "first_name"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	check(t, db.Select("*").Where("id = ?", control.ID).Get(&test))
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestSelectNoColumn(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	control := T{
		ID:   5,
		Name: "foo",
	}
	rows := sqlmock.NewRows([]string{"id", "unexpected"})
	rows.AddRow(control.ID, control.Name)
	mock.ExpectQuery(`^select \* from t where id = \? limit 1$`).WithArgs(control.ID).WillReturnRows(rows)
	var test T
	if err := db.Select("*").Where("id = ?", control.ID).Get(&test); err == nil {
		t.Fatalf("expected err")
	}
	check(t, mock.ExpectationsWereMet())
}
