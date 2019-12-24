package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// SizeOfFunc .
type SizeOfFunc func(reflect.StructField) uint64

// ErrNotFound .
var ErrNotFound = errors.New("no result found")

// Register .
func Register(structs ...interface{}) error {
	for _, s := range structs {
		if err := register(s); err != nil {
			return err
		}
	}
	return nil
}

type rowType struct {
	field string
	typ   string
	null  string
	key   string
	def   NullString
	extra NullString
}

func register(s interface{}) error {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("you can only register structs, %s found", reflect.TypeOf(s).Kind())
	}
	m := new(model)
	m.typ = typ
	m.name = m.typ.Name()
	m.table = toSnakeCase(m.name)
	m.primaryFieldIndex = -1
	for i := 0; i < m.typ.NumField(); i++ {
		f := m.typ.Field(i)
		tag, ok := f.Tag.Lookup("gosql")
		if ok && tag == "-" {
			continue
		}
		if ok && tag == "primary" {
			m.primaryFieldIndex = i
		}
		m.fields = append(m.fields, toSnakeCase(f.Name))
	}
	if err := m.mustBeValid(); err != nil {
		return err
	}
	m.fieldCount = len(m.fields)
	models[m.name] = m
	return nil
}

// Conn returns a reference to DB.
func Conn(db *sql.DB) *DB {
	return &DB{db}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// CheckSchema checks your models agains the database schema. This is
// only compatible with MySQL.
func CheckSchema(db *DB, sizer SizeOfFunc) error {
	uint64TypeRegexp, err := regexp.Compile(`(?i)^bigint\([0-9]+\) unsigned$`)
	if err != nil {
		return err
	}
	int64TypeRegexp, err := regexp.Compile(`(?i)^bigint\([0-9]+\)$`)
	if err != nil {
		return err
	}
	uintTypeRegexp, err := regexp.Compile(`(?i)^int\([0-9]+\) unsigned$`)
	if err != nil {
		return err
	}
	intTypeRegexp, err := regexp.Compile(`(?i)^int\([0-9]+\)$`)
	if err != nil {
		return err
	}
	boolTypeRegexp, err := regexp.Compile(`(?i)^tinyint\([0-9]+\)$`)
	if err != nil {
		return err
	}
	timeTypeRegexp, err := regexp.Compile(`(?i)^datetime$`)
	if err != nil {
		return err
	}
	for _, m := range models {
		rows, err := db.Query(fmt.Sprintf("describe %s", m.table))
		if err != nil {
			return err
		}
		defer rows.Close()
		var rowTypes []*rowType
		for rows.Next() {
			rt := new(rowType)
			if err := rows.Scan(&rt.field, &rt.typ, &rt.null, &rt.key, &rt.def, &rt.extra); err != nil {
				return err
			}
			rowTypes = append(rowTypes, rt)
		}
		for _, rt := range rowTypes {
			fieldIndex := -1
			for i := 0; i < m.typ.NumField(); i++ {
				f := m.typ.Field(i)
				if toSnakeCase(f.Name) == rt.field {
					fieldIndex = i
				}
			}
			if fieldIndex < 0 {
				return fmt.Errorf("model %s expected to have field %s", m.name, rt.field)
			}
			f := m.typ.Field(fieldIndex)
			t, ok := f.Tag.Lookup("gosql")
			if ok {
				if t == "-" {
					continue
				}
				if t == "primary" {
					if rt.key != "PRI" {
						return fmt.Errorf("expected table column %s.%s to have key 'PRI'", m.table, rt.field)
					}
				}
			}
			if !ok || t != "primary" {
				if rt.key == "PRI" {
					return fmt.Errorf("expected table column %s.%s not to have key 'PRI'", m.table, rt.field)
				}
			}
			switch f.Type.Kind() {
			case reflect.String:
				size := sizer(f)
				if size == 0 {
					return fmt.Errorf("expected model %s.%s to have a max", m.name, f.Name)
				}
				if strings.HasPrefix(rt.typ, "varchar") {
					prefix := fmt.Sprintf("varchar(%d)", size)
					if !strings.HasPrefix(rt.typ, prefix) {
						return fmt.Errorf("expected table column %s.%s to have type %s", m.table, rt.field, prefix)
					}
				} else if rt.typ == "tinytext" {
					if size != 256 {
						return fmt.Errorf("expected model %s.%s to have max %d", m.name, f.Name, 256)
					}
				} else if rt.typ == "text" {
					if size != 65535 {
						return fmt.Errorf("expected model %s.%s to have max %d", m.name, f.Name, 65535)
					}
				} else if rt.typ == "mediumtext" {
					if size != 16777215 {
						return fmt.Errorf("expected model %s.%s to have max %d", m.name, f.Name, 16777215)
					}
				} else if rt.typ == "longtext" {
					if size != 4294967295 {
						return fmt.Errorf("expected model field %s.%s to have max %d", m.name, f.Name, 4294967295)
					}
				} else {
					return fmt.Errorf("expected table column %s.%s to have type varchar, tinytext, text, mediumtext, or longtext", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Uint64:
				if !uint64TypeRegexp.MatchString(rt.typ) {
					return fmt.Errorf("expected table column %s.%s to have type bigint and be unsigned", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Int64:
				if !int64TypeRegexp.MatchString(rt.typ) {
					return fmt.Errorf("expected table column %s.%s to have type bigint and not be unsigned", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Uint:
				if !uintTypeRegexp.MatchString(rt.typ) {
					return fmt.Errorf("expected table column %s.%s to have type int and be unsigned", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Int:
				if !intTypeRegexp.MatchString(rt.typ) {
					return fmt.Errorf("expected table column %s.%s to have type int and not be unsigned", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Bool:
				if !boolTypeRegexp.MatchString(rt.typ) {
					return fmt.Errorf("expected table column %s.%s to have type tinyint and not be unsigned", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Float32:
				if rt.typ != "float" {
					return fmt.Errorf("expected table column %s.%s to have type float", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Float64:
				if rt.typ != "double" {
					return fmt.Errorf("expected table column %s.%s to have type double", m.table, rt.field)
				}
				if rt.null == "YES" {
					return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
				}
			case reflect.Slice:
				switch f.Type.Elem().Name() {
				case "uint8":
					if rt.typ != "blob" {
						return fmt.Errorf("expected table column %s.%s to have type blob", m.table, rt.field)
					}
					if rt.null == "YES" {
						return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
					}
				default:
					log.Printf("warning: skipping schema check for %s.%s\n", m.name, f.Name)
				}
			case reflect.Struct:
				switch f.Type.Name() {
				case "NullUint64":
					if !uint64TypeRegexp.MatchString(rt.typ) {
						return fmt.Errorf("expected table column %s.%s to have type bigint and be unsigned", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				case "NullInt64":
					if !int64TypeRegexp.MatchString(rt.typ) {
						return fmt.Errorf("expected table column %s.%s to have type bigint and not be unsigned", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				case "NullUint":
					if !uintTypeRegexp.MatchString(rt.typ) {
						return fmt.Errorf("expected table column %s.%s to have type int and be unsigned", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				case "NullInt":
					if !intTypeRegexp.MatchString(rt.typ) {
						return fmt.Errorf("expected table column %s.%s to have type int and not be unsigned", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				case "NullBool":
					if !boolTypeRegexp.MatchString(rt.typ) {
						return fmt.Errorf("expected table column %s.%s to have type tinyint and not be unsigned", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				case "Time":
					if !timeTypeRegexp.MatchString(rt.typ) {
						return fmt.Errorf("expected table column %s.%s to have type datetime", m.table, rt.field)
					}
					if rt.null != "NO" {
						return fmt.Errorf("expected table column %s.%s not to be nullable", m.table, rt.field)
					}
				case "NullTime":
					if !timeTypeRegexp.MatchString(rt.typ) {
						return fmt.Errorf("expected table column %s.%s to have type datetime", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				case "NullFloat64":
					if rt.typ != "double" {
						return fmt.Errorf("expected table column %s.%s to have type double", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				case "NullFloat32":
					if rt.typ != "float" {
						return fmt.Errorf("expected table column %s.%s to have type float", m.table, rt.field)
					}
					if rt.null != "YES" {
						return fmt.Errorf("expected table column %s.%s to be nullable", m.table, rt.field)
					}
				default:
					log.Printf("warning: skipping schema check for %s.%s\n", m.name, f.Name)
				}
			default:
				log.Printf("warning: skipping schema check for %s.%s\n", m.name, f.Name)
			}
		}
	}
	return nil
}
