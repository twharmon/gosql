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

	if _, err := DB.Insert(user); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertErrors(t *testing.T) {
	_, err := DB.Insert(User{})
	assertErr(t, "should return error if struct and not pointer to struct", err)

	_, err2 := DB.Insert(User{})
	assertErr(t, "should return error if pointer to non struct", err2)

	type Post struct {
		ID    int64
		Title string
	}
	_, err3 := DB.Insert(User{})
	assertErr(t, "should return error if struct not registered", err3)
}

func TestExec(t *testing.T) {
	mock.ExpectExec(`^update user set email = \? where id = \?$`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := DB.Exec("update user set email = ? where id = ?", 1); err != nil {
		t.Errorf("error was not expected while executing: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
