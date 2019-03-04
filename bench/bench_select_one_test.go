package bench

import (
	"log"
	"testing"
)

func BenchmarkSelectOneGOSQL(b *testing.B) {
	query := dbGosql.
		Query().
		Select("*").
		Where("id = ?", 163387)
	for i := 0; i < b.N; i++ {
		var user User
		err := query.To(&user)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkSelectOneGORM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var user User
		dbGorm.Where("id = ?", 163387).First(&user)
	}
}

func BenchmarkSelectOnePlain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := plainSelect("id = ?", 163387)
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
