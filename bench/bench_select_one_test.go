package bench

import (
	"log"
	"testing"
)

func BenchmarkSelectOneGOSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := User{}
		err := dbGosql.
			Select("*").
			Where("id = ?", 1).
			To(&user)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkSelectOneGORM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := User{}
		dbGorm.Where("id = ?", 1).First(&user)
	}
}

func BenchmarkSelectOnePlain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := plainSelect("id = ?", 1)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func plainSelect(where string, id int) (*User, error) {
	user := new(User)
	err := dbPlain.QueryRow("select * from user where "+where, id).Scan(
		&user.ID,
		&user.Role,
		&user.Email,
		&user.Active,
	)
	return user, err
}
