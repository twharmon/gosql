package bench

import (
	"log"
	"testing"
)

func BenchmarkInsertGOSQL(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = dbGosql.Insert(&User{
			Role:   "asdf",
			Email:  "fdsa",
			Active: true,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkInsertGORM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dbGorm.Create(&User{
			Role:   "asdf",
			Email:  "fdsa",
			Active: true,
		})
	}
}

func BenchmarkInsertPlain(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = plainInsert(&User{
			Role:   "asdf",
			Email:  "fdsa",
			Active: true,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func plainInsert(u *User) error {
	res, err := dbPlain.Exec(
		"insert into user (role, email, active) values (?, ?, ?)",
		u.Role,
		u.Email,
		u.Active,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}
