package gosql_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestBuilderSelectOneAllFields(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	rows.AddRow(control.ID, control.Role, control.Email, control.Active)
	mock.ExpectQuery(`^select \* from user where id = \?$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DB.Query().Select("*").Where("id = ?", 1).To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectOneSomeFields(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"email", "active"})
	rows.AddRow(control.Email, control.Active)
	mock.ExpectQuery(`^select email, active from user where id = \?$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DB.Query().Select("email", "active").Where("id = ?", 1).To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	test.ID = control.ID
	test.Role = control.Role
	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectOneOneField(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"email"})
	rows.AddRow(control.Email)
	mock.ExpectQuery(`^select email from user where id = \?$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DB.Query().Select("email").Where("id = ?", 1).To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	test.ID = control.ID
	test.ID = control.ID
	test.Role = control.Role
	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectOneJoin(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	rows.AddRow(control.ID, control.Role, control.Email, control.Active)
	mock.ExpectQuery(`^select user\.\* from user join post on post.user_id = user.id where id = \?$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DB.Query().Select("*").Where("id = ?", 1).Join("post on post.user_id = user.id").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyAllFields(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Role, c.Email, c.Active)
	}
	mock.ExpectQuery(`^select \* from user$`).WillReturnRows(rows)

	var test []*User
	if err := DB.Query().Select("*").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}
func TestBuilderSelectManySomeFields(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"email", "active"})
	for i, _ := range control {
		rows.AddRow(control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select email, active from user$`).WillReturnRows(rows)

	var test []*User
	if err := DB.Query().Select("email", "active").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i, _ := range control {
		test[i].ID = control[i].ID
		test[i].Role = control[i].Role
	}
	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}
func TestBuilderSelectManyOneField(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"email"})
	for i, _ := range control {
		rows.AddRow(control[i].Email)
	}
	mock.ExpectQuery(`^select email from user$`).WillReturnRows(rows)

	var test []*User
	if err := DB.Query().Select("email").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i, _ := range control {
		test[i].ID = control[i].ID
		test[i].Role = control[i].Role
		test[i].Active = control[i].Active
	}
	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyLimit(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	for i, _ := range control {
		rows.AddRow(control[i].ID, control[i].Role, control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select \* from user where id = \? limit 5$`).WithArgs(1).WillReturnRows(rows)

	var test []*User
	if err := DB.Query().Select("*").Where("id = ?", 1).Limit(5).To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyOrderBy(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	for i, _ := range control {
		rows.AddRow(control[i].ID, control[i].Role, control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select \* from user where id = \? order by email$`).WithArgs(1).WillReturnRows(rows)

	var test []*User
	if err := DB.Query().Select("*").Where("id = ?", 1).OrderBy("email").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyJoinSomeFields(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"email", "active"})
	for i, _ := range control {
		rows.AddRow(control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select user\.email, user\.active from user join post on post.user_id = user.id$`).WillReturnRows(rows)

	var test []*User
	if err := DB.Query().Select("email", "active").Join("post on post.user_id = user.id").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i, _ := range control {
		test[i].ID = control[i].ID
		test[i].Role = control[i].Role
	}
	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderDelete(t *testing.T) {
	mock.ExpectExec(`^delete from user where email = \?$`).WithArgs("test@example.com").WillReturnResult(sqlmock.NewResult(1, 1))

	if err := DB.Query().Delete("user").Where("email = ?", "test@example.com").Exec(); err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
