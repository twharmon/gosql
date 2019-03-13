package gosql_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestExec(t *testing.T) {
	mock.ExpectExec(`^update user set email = \? where id = \?$`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := DB.Exec("update user set email = ? where id = ?", 1); err != nil {
		t.Errorf("error was not expected while executing: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
