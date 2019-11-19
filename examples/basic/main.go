package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/twharmon/gosql"
)

// User .
type User struct {
	Email    string `json:"email"`
	UserID   int    `json:"id" gosql:"primary"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
	IsActive bool   `json:"isActive"`
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

	gosql.Register(User{})
	db = gosql.Conn(database)
}

func main() {
	var user User
	if err := db.Select("*").To(&user); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("user:", user)

	newUser := User{
		Email:    "asdf@example.com",
		Password: "asdf",
		IsAdmin:  true,
		IsActive: true,
		UserID:   99999999999,
	}
	res, err := db.Insert(&newUser)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.RowsAffected())
	fmt.Println(res.LastInsertId())
}
