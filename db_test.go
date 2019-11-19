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

func TestUpdate(t *testing.T) {
	user := &User{
		ID:     5,
		Role:   "admin",
		Email:  "test@example.com",
		Active: true,
	}

	mock.ExpectExec(`^update user set role = \?, email = \?, active = \? where id = \?$`).WithArgs(user.Role, user.Email, user.Active, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	q := DB.Update("user")
	q.Set("role = ?", user.Role)
	q.Set("email = ?", user.Email)
	q.Set("active = ?", user.Active)
	q.Where("id = ?", user.ID)

	if _, err := q.Exec(); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	mock.ExpectExec(`^delete from user where id = \?$`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	q := DB.Delete("user")
	q.Where("id = ?", 1)

	if _, err := q.Exec(); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCount(t *testing.T) {
	rows := sqlmock.NewRows([]string{"count(*)"})
	rows.AddRow(10)
	mock.ExpectQuery(`^select count\(\*\) from user where id > \?$`).WithArgs(10).WillReturnRows(rows)

	q := DB.Count("user", "*")
	q.Where("id > ?", 10)

	if _, err := q.Exec(); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertPrimaryLast(t *testing.T) {
	user := &UserPrimaryLast{
		Role:   "admin",
		Email:  "test@example.com",
		Active: true,
	}

	mock.ExpectExec(`^insert into user_primary_last \(role, email, active\) values \(\?, \?, \?\)$`).WithArgs(user.Role, user.Email, user.Active).WillReturnResult(sqlmock.NewResult(1, 1))

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
