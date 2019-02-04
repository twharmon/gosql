package bench

import (
	"log"
	"strconv"
	"testing"
)

func BenchmarkSelectManyGOSQL(b *testing.B) {
	var users []*User
	for i := 0; i < b.N; i++ {
		err := dbGosql.
			Select("*").
			Where("id < ?", 100).
			Limit(1000).
			To(&users)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkSelectManyGORM(b *testing.B) {
	var users []*User
	for i := 0; i < b.N; i++ {
		dbGorm.
			Where("id < ?", 100).
			Limit(1000).
			Find(&users)
	}
}

func BenchmarkSelectManyPlain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := plainSelectMany("id < ?", 100, 1000)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func plainSelectMany(where string, id int, limit int) ([]*User, error) {
	rows, err := dbPlain.Query("select * from user where "+where+" limit "+strconv.Itoa(limit), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := new(User)
		if err := rows.Scan(&user.ID, &user.Role, &user.Email, &user.Active); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
