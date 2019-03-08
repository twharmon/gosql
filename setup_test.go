package gosql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/twharmon/gofake"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/twharmon/gosql"
)

var DB *gosql.DB
var mock sqlmock.Sqlmock

// User contains user information
type User struct {
	ID     int64  `json:"id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

// ExpandedUser contains user information, and other fields not in db
type ExpandedUser struct {
	ID     int64  `json:"id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
	Friend *User  `json:"friend"`
}

func init() {
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	DB = gosql.Conn(db)
	gosql.Register(User{}, ExpandedUser{})
}

func assertSame(control *User, test *User) error {
	if control.ID != test.ID {
		return fmt.Errorf("test did not match control: control.ID %d, test.ID %d", control.ID, test.ID)
	}
	if control.Role != test.Role {
		return fmt.Errorf("test did not match control: control.Role %s, test.Role %s", control.Role, test.Role)
	}
	if control.Email != test.Email {
		return fmt.Errorf("test did not match control: control.Email %s, test.Email %s", control.Email, test.Email)
	}
	if control.Active != test.Active {
		return fmt.Errorf("test did not match control: control.Active %t, test.Active %t", control.Active, test.Active)
	}
	return nil
}

func assertExpandedSame(control *ExpandedUser, test *ExpandedUser) error {
	if control.ID != test.ID {
		return fmt.Errorf("test did not match control: control.ID %d, test.ID %d", control.ID, test.ID)
	}
	if control.Role != test.Role {
		return fmt.Errorf("test did not match control: control.Role %s, test.Role %s", control.Role, test.Role)
	}
	if control.Email != test.Email {
		return fmt.Errorf("test did not match control: control.Email %s, test.Email %s", control.Email, test.Email)
	}
	if control.Active != test.Active {
		return fmt.Errorf("test did not match control: control.Active %t, test.Active %t", control.Active, test.Active)
	}
	return nil
}

func assertSameSlice(control []*User, test []*User) error {
	if len(control) != len(test) {
		return fmt.Errorf("control hand length %d, but test had length %d", len(control), len(test))
	}
	for i, _ := range control {
		if err := assertSame(control[i], test[i]); err != nil {
			return err
		}
	}
	return nil
}

func assertPanic(t *testing.T, desc string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("assert panic: %s (no panic)", desc)
		}
	}()
	f()
}

func assertErr(t *testing.T, desc string, err error) {
	if err == nil {
		t.Errorf("assert error: %s (nil error)", desc)
	}
}

func makeUser() *User {
	return &User{
		ID:     int64(gofake.Int(100)),
		Role:   gofake.Word(),
		Email:  gofake.Email(),
		Active: gofake.Int(1) == 1,
	}
}

func makeExpandedUser() *ExpandedUser {
	return &ExpandedUser{
		ID:     int64(gofake.Int(100)),
		Role:   gofake.Word(),
		Email:  gofake.Email(),
		Active: gofake.Int(1) == 1,
		Friend: makeUser(),
	}
}

func makeUserSlice(size int) []*User {
	users := make([]*User, size)
	for i := 0; i < size; i++ {
		users[i] = makeUser()
	}
	return users
}
