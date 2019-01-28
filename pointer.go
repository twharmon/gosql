package gosql

import (
	"reflect"
)

func isPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr
}
