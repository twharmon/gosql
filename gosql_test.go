package gosql

import (
	"testing"
)

func TestRegister(t *testing.T) {
	m := make(map[string]string)
	assertErr(t, "should return error if non struct passed to Register", Register(m))

	type T struct {
		ID uint64 `gosql:"primary"`
	}
	assertErr(t, "should return error if struct registered twice", Register(T{}, T{}))
}

type User2 struct {
	ID    int `gosql:"primary"`
	Name  string
	Posts []*Post2
}

type Post2 struct {
	ID     int `gosql:"primary"`
	Title  string
	Author *User2
}

type User3 struct {
	ID    int `gosql:"primary"`
	Name  string
	Posts []*Post3
}

type Post3 struct {
	ID      int `gosql:"primary"`
	Title   string
	Authors []*User3
}

func TestOneToManyRegistration(t *testing.T) {
	// if err := Register(User2{}, Post2{}); err != nil {
	// 	t.Fatal("unexpected err", err)
	// }
	// otms := models["User2"].oneToManys
	// if otms["post2"] == nil {
	// 	t.Errorf("one to many failed to register")
	// }
	// mtos := models["Post2"].manyToOnes
	// if mtos["user2"] == nil {
	// 	t.Errorf("many to one failed to register")
	// }
}

func TestManyToManyRegistration(t *testing.T) {
	// if err := Register(User3{}, Post3{}); err != nil {
	// 	t.Fatal("unexpected err", err)
	// }
	// mtm1s := models["User3"].manyToManys
	// if mtm1s["post3"] == nil {
	// 	t.Errorf("many to many failed to register")
	// }
	// mtm2s := models["Post3"].manyToManys
	// if mtm2s["user3"] == nil {
	// 	t.Errorf("many to many failed to register")
	// }
}
