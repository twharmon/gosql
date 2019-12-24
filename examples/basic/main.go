package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/twharmon/gosql"
)

// User .
type User struct {
	Email    string `json:"email" size:"255"`
	ID       uint   `json:"id" gosql:"primary"`
	Password string `json:"password" size:"100"`
	IsAdmin  bool   `json:"isAdmin"`
	IsActive bool   `json:"isActive"`
	F        gosql.NullFloat32
	D        float64
}

var db *gosql.DB

func init() {
	database, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@%s/%s",
		"root",
		"",
		"",
		"test_db",
	))
	if err != nil {
		log.Fatalln(err)
	}

	if err := gosql.Register(User{}); err != nil {
		log.Fatalln(err)
	}
	db = gosql.Conn(database)
	if err := gosql.CheckSchema(db, sizeOf); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var user User
	if err := db.Select("*").To(&user); err != nil {
		if !errors.Is(err, gosql.ErrNotFound) {
			log.Fatalln(err)
		}
	}
	fmt.Println("user:", user)

	newUser := User{
		Email:    "asdf@example.com",
		Password: "asdf",
		IsAdmin:  true,
		IsActive: true,
	}
	res, err := db.Insert(&newUser)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.RowsAffected())
	fmt.Println(res.LastInsertId())
}

func sizeOf(f reflect.StructField) uint64 {
	tag, ok := f.Tag.Lookup("size")
	if !ok {
		return 0
	}
	size, err := strconv.ParseUint(tag, 10, 64)
	if err != nil {
		panic(err)
	}
	return size
}
