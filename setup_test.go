package gosql_test

import (
	"database/sql"
	"fmt"

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

func init() {
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	DB = gosql.ConnDB(db)
	gosql.Register(User{})
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

func makeUser() *User {
	return &User{
		ID:     int64(gofake.Int(100)),
		Role:   gofake.Word(),
		Email:  gofake.Email(),
		Active: gofake.Int(1) == 1,
	}
}

func makeUserSlice(size int) []*User {
	users := make([]*User, size)
	for i := 0; i < size; i++ {
		users[i] = makeUser()
	}
	return users
}
