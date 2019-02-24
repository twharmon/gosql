package bench

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/twharmon/gosql"
)

// User .
type User struct {
	ID     int64  `json:"id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

var dbGosql *gosql.DB
var dbPlain *sql.DB
var dbGorm *gorm.DB

func init() {
	var err error
	dbPlain, err = sql.Open("mysql", fmt.Sprintf(
		"%s:%s@%s/%s",
		"root",
		"",
		"",
		"test_db",
	))
	if err != nil {
		log.Fatalln(err)
	}

	dbGorm, err = gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@%s/%s",
		"root",
		"",
		"",
		"test_db",
	))
	if err != nil {
		log.Fatalln(err)
	}
	dbGorm.SingularTable(true)

	gosql.Register(&User{})
	dbGosql, err = gosql.Conn(
		"root",
		"",
		"",
		"test_db",
	)
	if err != nil {
		log.Fatalln(err)
	}
}
