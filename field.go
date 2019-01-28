package gosql

import (
	"reflect"
	"strings"
)

var fieldKinds = []reflect.Kind{
	reflect.String,
	reflect.Bool,
	reflect.Uint,
	reflect.Uint64,
	reflect.Int,
	reflect.Int64,
	reflect.Float32,
	reflect.Float64,
}

func isField(s reflect.StructField) bool {
	k := s.Type.Kind()
	for i := 0; i < len(fieldKinds); i++ {
		if fieldKinds[i] == k {
			return true
		}
	}
	return false
}

func isOneToMany(src *model, s reflect.StructField) bool {
	if s.Type.Kind() != reflect.Slice {
		return false
	}
	tarFullName := s.Type.String()
	tarName := tarFullName[strings.Index(tarFullName, ".")+1:]

	tar := models[tarName]
	if tar == nil {
		return false
	}

	is := false
	for i := 0; i < tar.typ.NumField(); i++ {
		if strings.HasSuffix(tar.typ.Field(i).Type.String(), src.name) {
			is = true
		}
	}
	return is
}

func isManyToOne(src *model, s reflect.StructField) bool {
	if s.Type.Kind() != reflect.Struct {
		return false
	}
	tarFullName := s.Type.String()
	tarName := tarFullName[strings.Index(tarFullName, ".")+1:]
	tar := models[tarName]
	if tar == nil {
		return false
	}

	is := false
	for i := 0; i < tar.typ.NumField(); i++ {
		if strings.HasSuffix(tar.typ.Field(i).Type.String(), src.name) {
			is = true
		}
	}
	return is
}

func isOneToOne(src *model, s reflect.StructField) bool {
	return false
}

func isManyToMany(src *model, s reflect.StructField) bool {
	return false
}
