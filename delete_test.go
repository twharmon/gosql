package gosql_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestDelete(t *testing.T) {
	user := &User{ID: 5}

	mock.ExpectExec(`^delete from user where id = \?$`).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := DB.Delete(user); err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
