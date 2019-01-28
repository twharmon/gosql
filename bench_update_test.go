package gosql

import (
	"log"
	"testing"
)

func BenchmarkUpdateGOSQL(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = dbGosql.Update(&User{
			ID:     2,
			Role:   "asdf",
			Email:  "fdsa",
			Active: true,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkUpdateGORM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dbGorm.Save(&User{
			ID:     2,
			Role:   "asdf",
			Email:  "fdsa",
			Active: true,
		})
	}
}

func BenchmarkUpdatePlain(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = plainUpdate(&User{
			ID:     2,
			Role:   "asdf",
			Email:  "fdsa",
			Active: true,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func plainUpdate(u *User) error {
	_, err := dbPlain.Exec(
		"update user set role = ?, email = ?, active = ? where id = ? limit 1",
		u.Role,
		u.Email,
		u.Active,
		u.ID,
	)
	return err
}
