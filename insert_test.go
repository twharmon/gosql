package gosql_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	user := &User{
		Role:   "admin",
		Email:  "test@example.com",
		Active: true,
	}

	mock.ExpectExec(`^insert into user \(role, email, active\) values \(\?, \?, \?\)$`).WithArgs(user.Role, user.Email, user.Active).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := DB.Insert(user); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertErrors(t *testing.T) {
	assertErr(t, "should return error if struct and not pointer to struct", DB.Insert(User{}))

	testMap := make(map[string]interface{})
	assertErr(t, "should return error if pointer to non struct", DB.Insert(&testMap))

	type Post struct {
		ID    int64
		Title string
	}
	assertErr(t, "should return error if struct not registered", DB.Insert(&Post{}))
}
