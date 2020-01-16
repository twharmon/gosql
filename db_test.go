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

type User struct {
	ID   int `gosql:"primary"`
	Name string
}

func init() {
	if err := gosql.Register(User{}); err != nil {
		panic(err)
	}
}

func TestDelete(t *testing.T) {
	type DeleteModel struct {
		ID int `gosql:"primary"`
	}
	check(t, gosql.Register(DeleteModel{}))
	deleteModel := DeleteModel{5}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^delete from delete_model where id = \?$`).WithArgs(deleteModel.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Delete(&deleteModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	type UpdateModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(UpdateModel{}))
	updateModel := UpdateModel{5, "foo"}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^update update_model set name = \? where id = \?$`).WithArgs(updateModel.Name, updateModel.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Update(&updateModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsert(t *testing.T) {
	type InsertModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(InsertModel{}))
	insertModel := InsertModel{Name: "foo"}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^insert into insert_model \(name\) values \(\?\)$`).WithArgs(insertModel.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&insertModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsertWithPrimary(t *testing.T) {
	type InsertWithPrimaryModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(InsertWithPrimaryModel{}))
	insertModelWithPrimary := InsertWithPrimaryModel{5, "foo"}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^insert into insert_with_primary_model \(id, name\) values \(\?, \?\)$`).WithArgs(insertModelWithPrimary.ID, insertModelWithPrimary.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&insertModelWithPrimary)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func ExampleDB_Insert() {
	os.Remove("/tmp/foo.db")
	sqliteDB, _ := sql.Open("sqlite3", "/tmp/foo.db")
	sqliteDB.Exec("create table user (id integer not null primary key, name text); delete from user")
	db := gosql.Conn(sqliteDB)
	type User struct {
		ID   int `gosql:"primary"`
		Name string
	}
	gosql.Register(User{})
	db.Insert(&User{Name: "Gopher"})
	var user User
	db.Select("*").To(&user)
	fmt.Println(user.Name)
	// Output: Gopher
}

func ExampleDB_Update() {
	os.Remove("/tmp/foo.db")
	sqliteDB, _ := sql.Open("sqlite3", "/tmp/foo.db")
	sqliteDB.Exec("create table user (id integer not null primary key, name text); delete from user")
	db := gosql.Conn(sqliteDB)
	type User struct {
		ID   int `gosql:"primary"`
		Name string
	}
	gosql.Register(User{})
	user := User{ID: 5, Name: "Gopher"}
	db.Insert(&user)
	user.Name = "Gofer"
	db.Update(&user)
	var foo User
	db.Select("*").To(&foo)
	fmt.Println(foo.Name)
	// Output: Gofer
}

func ExampleDB_Delete() {
	os.Remove("/tmp/foo.db")
	sqliteDB, _ := sql.Open("sqlite3", "/tmp/foo.db")
	sqliteDB.Exec("create table user (id integer not null primary key, name text); delete from user")
	db := gosql.Conn(sqliteDB)
	type User struct {
		ID   int `gosql:"primary"`
		Name string
	}
	gosql.Register(User{})
	user := User{ID: 5, Name: "Gopher"}
	db.Insert(&user)
	db.Delete(&user)
	var foo User
	err := db.Select("*").To(&foo)
	fmt.Println(err)
	// Output: no result found
}

func BenchmarkInsert(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	user := User{Name: "Gopher"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.Insert(&user)
		check(b, err)
	}
}

func BenchmarkUpdate(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
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
	user := User{ID: 5, Name: "Gopher"}
	_, err := db.Insert(&user)
	check(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		check(b, db.Select("*").To(&user))
	}
}

func BenchmarkSelectMany(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	user := User{Name: "Gopher"}
	for i := 0; i < 100; i++ {
		_, err := db.Insert(&user)
		check(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		check(b, db.Select("*").Limit(100).To(&users))
	}
}

func BenchmarkSelectManyPtrs(b *testing.B) {
	db := getSQLiteDB(b, "create table user (id integer not null primary key, name text); delete from user")
	user := User{Name: "Gopher"}
	for i := 0; i < 100; i++ {
		_, err := db.Insert(&user)
		check(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []*User
		check(b, db.Select("*").Limit(100).To(&users))
	}
}
