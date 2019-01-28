package gosql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

// User .
type User struct {
	ID     uint   `json:"id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

var dbGosql *DB
var dbPlain *sql.DB
var dbGorm *gorm.DB

func init() {
	var err error
	dbPlain, err = sql.Open("mysql", fmt.Sprintf(
		"%s:%s@%s/%s",
		"root",
		"",
		"",
		"solasites",
	))
	if err != nil {
		log.Fatalln(err)
	}

	dbGorm, err = gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@%s/%s",
		"root",
		"",
		"",
		"solasites",
	))
	if err != nil {
		log.Fatalln(err)
	}
	dbGorm.SingularTable(true)

	err = MustPrepare(&User{})
	if err != nil {
		log.Fatalln(err)
	}
	dbGosql, err = Conn(
		"root",
		"",
		"",
		"solasites",
	)
	if err != nil {
		log.Fatalln(err)
	}
}
