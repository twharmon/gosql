package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestManualUpdate(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	val1 := 1
	val2 := 2
	mock.ExpectExec(`^update t set val1 = \?, val2 = \?$`).WithArgs(val1, val2).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.ManualUpdate("t").
		Set("val1 = ?", val1).
		Set("val2 = ?", val2).
		Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestManualUpdateWhere(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	val1 := 1
	val2 := 2
	val3 := 3
	val4 := 4
	mock.ExpectExec(`^update t set val1 = \? where val2 = \? and val3 = \? or val4 = \?$`).WithArgs(val1, val2, val3, val4).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.ManualUpdate("t").
		Set("val1 = ?", val1).
		Where("val2 = ?", val2).
		Where("val3 = ?", val3).
		OrWhere("val4 = ?", val4).
		Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestManualUpdateJoin(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	val := 1
	mock.ExpectExec(`^update t join a on a.id = t.a_id set val = \?$`).WithArgs(val).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.ManualUpdate("t").
		Set("val = ?", val).
		Join("a on a.id = t.a_id").
		Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}
