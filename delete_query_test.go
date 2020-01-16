package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestManualDelete(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	val := 1
	mock.ExpectExec(`^delete from t where val = \?$`).WithArgs(val).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.ManualDelete("t").Where("val = ?", val).Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestManualDeleteWhere(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	val1 := 1
	val2 := 2
	val3 := 3
	mock.ExpectExec(`^delete from t where val1 = \? and val2 = \? or val3 = \?$`).WithArgs(val1, val2, val3).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.ManualDelete("t").
		Where("val1 = ?", val1).
		Where("val2 = ?", val2).
		OrWhere("val3 = ?", val3).
		Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestManualDeleteJoin(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^delete from t join a on a.id = t.a_id$`).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.ManualDelete("t").Join("a on a.id = t.a_id").Exec()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}
