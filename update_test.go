package gosql_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestUpdate(t *testing.T) {
	user := &User{
		Role:   "admin",
		Email:  "test@example.com",
		Active: true,
	}

	mock.ExpectExec(`^update user set role = \?, email = \?, active = \? where id = \?$`).WithArgs(user.Role, user.Email, user.Active, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := DB.Update(user); err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
