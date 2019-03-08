package gosql_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestBuilderDelete(t *testing.T) {
	mock.ExpectExec(`^delete from user where email = \?$`).WithArgs("test@example.com").WillReturnResult(sqlmock.NewResult(1, 1))

	if err := DB.Table("user").Where("email = ?", "test@example.com").Delete(); err != nil {
		t.Errorf("error was not expected while deleting: %v", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBuilderCount(t *testing.T) {
	mock.ExpectQuery(`^select count\(\*\) from user$`).WillReturnRows(mock.NewRows([]string{"count(*)"}).AddRow(10))

	count, err := DB.Table("user").Count()
	if err != nil {
		t.Errorf("error was not expected while counting: %v", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if count != 10 {
		t.Errorf("expected count of 10, got %d", count)
	}
}
