package main

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/twharmon/gosql"
)

// User .
type User struct {
	ID       uint64 `json:"id" gosql:"primary"`
	Email    string `json:"email" size:"255"`
	Password string `json:"password" size:"100"`
	IsAdmin  bool   `json:"isAdmin"`
	IsActive bool   `json:"isActive"`
	F        gosql.NullFloat32
	D        float64
	Posts    []*Post
}

// Post .
type Post struct {
	ID        uint64 `json:"id" gosql:"primary"`
	Title     string `size:"255"`
	Body      string `size:"65535"`
	Published bool
	Author    *User
}

var db *gosql.DB

// Scan implements the Scanner interface.
func (u *User) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch value.(type) {
	case int32:
		u.ID = uint64(value.(int32))
	case int:
		u.ID = uint64(value.(int))
	case int64:
		u.ID = uint64(value.(int64))
	// case int32:
	// 	u.ID = uint64(value.(int32))
	// case int32:
	// 	u.ID = uint64(value.(int32))
	case []byte:
		i, err := strconv.ParseUint(string(value.([]byte)), 10, 64)
		if err != nil {
			return err
		}
		u.ID = i
	default:
		return errors.New("switch failed")
	}

	return nil
}

// Value implements the driver Valuer interface.
func (u User) Value() (driver.Value, error) {
	return u.ID, nil
}

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

	if err := gosql.Register(User{}, Post{}); err != nil {
		log.Fatalln(err)
	}
	db = gosql.Conn(database)
	// if err := gosql.CheckSchema(db, sizeOf); err != nil {
	// 	log.Fatalln(err)
	// }
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
	res, err := db.Save(&newUser)
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

/*

err := models.DB.Insert(&user)
err := models.DB.Update(&user)
err := models.DB.Delete(&user)

models.DB.Select("*")
models.DB.ManualUpdate("user") ManualUpdate?
models.DB.ManualDelete("user") ManualDelete?

// err := models.DB.On(&post).Load("author").Select("id", "email").Exec()
// err := models.DB.On(&user).Load("posts").Select("id", "title").Exec()
// ---------------------------------------------------------------------
// err := models.DB.On(&post).Load("Author").Select("id", "email").Exec()
// err := models.DB.On(&user).Load("Posts").Select("id", "title").Exec()
// err := models.DB.On(&user).
// 	Insert()
// 	Update()
// 	Delete()
*/
