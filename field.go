package gosql

import (
	"reflect"
	"strings"
)

var fieldTypes = []string{
	"string",
	"bool",
	"uint",
	"uint64",
	"int",
	"int64",
	"float32",
	"float64",
	"[]byte",
	"gosql.NullInt64",
	"gosql.NullString",
	"gosql.NullFloat64",
	"gosql.NullBool",
	"gosql.NullTime",
	"time.Time",
}

func isField(s reflect.StructField) bool {
	t := s.Type.String()
	for i := 0; i < len(fieldTypes); i++ {
		firstLetter := string(s.Name[0])
		if fieldTypes[i] == t && firstLetter == strings.ToUpper(firstLetter) {
			return true
		}
	}
	return false
}
