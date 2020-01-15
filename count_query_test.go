package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCountStar(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	control := int64(15)
	rows := sqlmock.NewRows([]string{"count(*)"})
	rows.AddRow(control)
	mock.ExpectQuery(`^select count\(\*\) from t$`).WillReturnRows(rows)
	test, err := db.Count("t", "*").Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestCountStarWhere(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	control := int64(15)
	val1 := 1
	val2 := 2
	val3 := 3
	rows := sqlmock.NewRows([]string{"count(*)"})
	rows.AddRow(control)
	mock.ExpectQuery(`^select count\(\*\) from t where val1 > \? and val2 > \? or val3 > \?$`).WithArgs(val1, val2, val3).WillReturnRows(rows)
	test, err := db.Count("t", "*").
		Where("val1 > ?", val1).
		Where("val2 > ?", val2).
		OrWhere("val3 > ?", val3).
		Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}

func TestCountStarJoin(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	control := int64(15)
	rows := sqlmock.NewRows([]string{"count(*)"})
	rows.AddRow(control)
	mock.ExpectQuery(`^select count\(\*\) from t join a on a.id = t.a_id$`).WillReturnRows(rows)
	test, err := db.Count("t", "*").Join("a on a.id = t.a_id").Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
	equals(t, control, test)
}
