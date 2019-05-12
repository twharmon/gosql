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
