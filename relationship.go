package gosql

import (
	"reflect"
	"strings"
)

// User -> Post
type oneToMany struct {
	typStr string              // "Post"
	field  reflect.StructField // Posts []*Post
}

// Post -> User
type manyToOne struct {
	typStr string              // "User"
	field  reflect.StructField // Author *User
	column string              // "author_id"
}

// Post -> Tag
type manyToMany struct {
	typStr string              // "Tag"
	field  reflect.StructField // Tags []*Tag
	table  string              // "post_tag"
	column string              // "tag_id"
}

func (m *model) getOneToManyPairByType(typ reflect.Type) (*oneToMany, *manyToOne) {
	for _, otm := range m.oneToManys {
		if strings.HasSuffix(typ.String(), otm.typStr) {
			r := models[strings.Split(otm.typStr, ".")[1]]
			for _, mto := range r.manyToOnes {
				if strings.HasSuffix(mto.typStr, m.typ.String()) {
					return otm, mto
				}
			}
		}
	}
	return nil, nil
}

func (m *model) isOneToMany(s reflect.StructField) bool {
	if s.Type.Kind() != reflect.Slice {
		return false
	}
	tarFullName := s.Type.String() // []*models.Post | []*Post
	if !strings.HasPrefix(tarFullName, "[]*") {
		return false
	}
	tarName := tarFullName[strings.Index(tarFullName, ".")+1:]
	tar := models[tarName]
	if tar == nil {
		return false
	}
	is := false
	for i := 0; i < tar.typ.NumField(); i++ {
		typ := tar.typ.Field(i).Type.String()
		if strings.HasSuffix(typ, "."+m.name) && strings.HasPrefix(typ, "*") {
			is = true
		}
	}
	return is
}

func (m *model) isManyToOne(s reflect.StructField) bool {
	if s.Type.Kind() != reflect.Ptr {
		return false
	}
	tarFullName := s.Type.String()
	if strings.HasPrefix(tarFullName, "[]*") {
		return false
	}
	tarName := tarFullName[strings.Index(tarFullName, ".")+1:]
	tar := models[tarName]
	if tar == nil {
		return false
	}

	is := false
	for i := 0; i < tar.typ.NumField(); i++ {
		typ := tar.typ.Field(i).Type.String()
		if strings.HasSuffix(typ, "."+m.name) && strings.HasPrefix(typ, "[]*") {
			is = true
		}
	}
	return is
}

func (m *model) isManyToMany(s reflect.StructField) bool {
	if s.Type.Kind() != reflect.Slice {
		return false
	}
	tarFullName := s.Type.String()
	if !strings.HasPrefix(tarFullName, "[]*") {
		return false
	}
	tarName := tarFullName[strings.Index(tarFullName, ".")+1:]
	tar := models[tarName]
	if tar == nil {
		return false
	}

	is := false
	for i := 0; i < tar.typ.NumField(); i++ {
		typ := tar.typ.Field(i).Type.String() // []*models.Post
		if strings.HasSuffix(typ, "."+m.name) && strings.HasPrefix(typ, "[]*") {
			is = true
		}
	}
	return is
}
