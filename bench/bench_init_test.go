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
	ID     uint   `json:"id"`
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

	err = gosql.MustPrepare(&User{})
	if err != nil {
		log.Fatalln(err)
	}
	dbGosql, err = gosql.Conn(
		"root",
		"",
		"",
		"solasites",
	)
	if err != nil {
		log.Fatalln(err)
	}
}
