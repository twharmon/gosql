package gosql_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/twharmon/gosql"
)

func TestDelete(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID int `idx:"primary"`
	}
	deleteModel := T{5}
	mock.ExpectExec(`^delete from t where id = \?$`).WithArgs(deleteModel.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Delete(&deleteModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	updateModel := T{5, "foo"}
	mock.ExpectExec(`^update t set name = \? where id = \?$`).WithArgs(updateModel.Name, updateModel.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Update(&updateModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestUpdateThreeFields(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID    int `idx:"primary"`
		Name  string
		Email string
	}
	updateModel := T{5, "foo", "foo@example.com"}
	mock.ExpectExec(`^update t set name = \?, email = \? where id = \?$`).WithArgs(updateModel.Name, updateModel.Email, updateModel.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Update(&updateModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestUpdateThreeFieldsTwoPrimaries(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID    int `idx:"primary"`
		Name  string
		Email string `idx:"primary"`
	}
	updateModel := T{5, "foo", "foo@example.com"}
	mock.ExpectExec(`^update t set name = \? where id = \? and email = \?$`).WithArgs(updateModel.Name, updateModel.ID, updateModel.Email).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Update(&updateModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestBegin(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	_ = mock.ExpectBegin()
	_, err = db.Begin()
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsert(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	insertModel := T{Name: "foo"}
	mock.ExpectExec(`^insert into t \(name\) values \(\?\)$`).WithArgs(insertModel.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&insertModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsertWithPrimary(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int `idx:"primary"`
		Name string
	}
	insertModelWithPrimary := T{5, "foo"}
	mock.ExpectExec(`^insert into t \(id, name\) values \(\?, \?\)$`).WithArgs(insertModelWithPrimary.ID, insertModelWithPrimary.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&insertModelWithPrimary)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsertWithAllFieldsPrimary(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID   int    `idx:"primary"`
		Name string `idx:"primary"`
	}
	model := T{5, "foo"}
	mock.ExpectExec(`^insert into t \(id, name\) values \(\?, \?\)$`).WithArgs(model.ID, model.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&model)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsertWith1stAndLastFieldsPrimary(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID    int `idx:"primary"`
		Email string
		Name  string `idx:"primary"`
	}
	model := T{5, "", "foo"}
	mock.ExpectExec(`^insert into t \(id, email, name\) values \(\?, \?, \?\)$`).WithArgs(model.ID, model.Email, model.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&model)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsertWithOmittedField(t *testing.T) {
	db, mock, err := getMockDB()
	check(t, err)
	type T struct {
		ID    int `idx:"primary"`
		Email string
		Name  string `col:"-"`
	}
	model := T{5, "", "foo"}
	mock.ExpectExec(`^insert into t \(id, email\) values \(\?, \?\)$`).WithArgs(model.ID, model.Email).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&model)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func ExampleDB_Insert() {
	os.Remove("/tmp/foo.db")
	sqliteDB, _ := sql.Open("sqlite3", "/tmp/foo.db")
	sqliteDB.Exec("create table user (id integer not null primary key, name text); delete from user")
	db := gosql.New(sqliteDB)
	type User struct {
		ID   int `idx:"primary"`
		Name string
	}
	db.Insert(&User{Name: "Gopher"})
	var user User
	db.Select("*").Get(&user)
	fmt.Println(user.Name)
	// Output: Gopher
}

func ExampleDB_Update() {
	os.Remove("/tmp/foo.db")
	sqliteDB, _ := sql.Open("sqlite3", "/tmp/foo.db")
	sqliteDB.Exec("create table user (id integer not null primary key, name text, email text); delete from user")
	db := gosql.New(sqliteDB)
	type User struct {
		ID    int `idx:"primary"`
		Name  string
		Email string
	}
	user := User{ID: 5, Name: "Gopher", Email: "gopher@example.com"}
	db.Insert(&user)
	user.Name = "Gofer"
	user.Email = "gofer@example.com"
	db.Update(&user)
	var foo User
	db.Select("*").Get(&foo)
	fmt.Println(foo.Name, foo.Email)
	// Output: Gofer gofer@example.com
}

func ExampleDB_Delete() {
	os.Remove("/tmp/foo.db")
	sqliteDB, _ := sql.Open("sqlite3", "/tmp/foo.db")
	sqliteDB.Exec("create table user (id integer not null primary key, name text); delete from user")
	db := gosql.New(sqliteDB)
	type User struct {
		ID   int `idx:"primary"`
		Name string
	}
	user := User{ID: 5, Name: "Gopher"}
	db.Insert(&user)
	db.Delete(&user)
	var foo User
	err := db.Select("*").Get(&foo)
	fmt.Println(err)
	// Output: no result found
}

func BenchmarkInsert(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	type User struct {
		ID   int `idx:"primary"`
		Name string
	}
	user := User{Name: "Gopher"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.Insert(&user)
		check(b, err)
	}
}

func BenchmarkUpdate(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	type User struct {
		ID   int `idx:"primary"`
		Name string
	}
	user := User{Name: "Gopher"}
	_, err := db.Insert(&user)
	check(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.Update(&user)
		check(b, err)
	}
}

func BenchmarkSelect(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	type User struct {
		ID   int `idx:"primary"`
		Name string
	}
	user := User{ID: 5, Name: "Gopher"}
	_, err := db.Insert(&user)
	check(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		check(b, db.Select("*").Get(&user))
	}
}

func BenchmarkSelectMany(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	type User struct {
		ID   int `idx:"primary"`
		Name string
	}
	user := User{Name: "Gopher"}
	for i := 0; i < 100; i++ {
		_, err := db.Insert(&user)
		check(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		check(b, db.Select("*").Limit(100).Get(&users))
	}
}

func BenchmarkSelectManyPtrs(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	type User struct {
		ID   int `idx:"primary"`
		Name string
	}
	user := User{Name: "Gopher"}
	for i := 0; i < 100; i++ {
		_, err := db.Insert(&user)
		check(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []*User
		check(b, db.Select("*").Limit(100).Get(&users))
	}
}
