package gosql

import (
	"database/sql"
	"strings"
)

const (
	sqlActionDelete = iota
	sqlActionCount
	sqlActionUpdate
)

// TableQuery .
type TableQuery struct {
	db        *DB
	action    int
	table     string
	joins     []string
	wheres    []*where
	sets      []string
	whereArgs []interface{}
	setArgs   []interface{}
}

// Table .
func (db *DB) Table(table string) *TableQuery {
	tq := new(TableQuery)
	tq.db = db
	tq.table = table
	return tq
}

// Where .
func (tq *TableQuery) Where(condition string, args ...interface{}) *TableQuery {
	w := &where{
		conjunction: " and ",
		condition:   condition,
	}
	tq.wheres = append(tq.wheres, w)
	tq.whereArgs = append(tq.whereArgs, args...)
	return tq
}

// OrWhere .
func (tq *TableQuery) OrWhere(condition string, args ...interface{}) *TableQuery {
	w := &where{
		conjunction: " or ",
		condition:   condition,
	}
	tq.wheres = append(tq.wheres, w)
	tq.whereArgs = append(tq.whereArgs, args...)
	return tq
}

// Set .
func (tq *TableQuery) Set(set string, arg interface{}) *TableQuery {
	tq.sets = append(tq.sets, set)
	tq.setArgs = append(tq.setArgs, arg)
	return tq
}

// Join .
func (tq *TableQuery) Join(join string) *TableQuery {
	tq.joins = append(tq.joins, join)
	return tq
}

// Delete .
func (tq *TableQuery) Delete() (sql.Result, error) {
	tq.action = sqlActionDelete
	args := tq.setArgs
	args = append(args, tq.whereArgs...)
	return tq.db.db.Exec(tq.string(), args...)
}

// Update .
func (tq *TableQuery) Update() (sql.Result, error) {
	tq.action = sqlActionUpdate
	args := tq.setArgs
	args = append(args, tq.whereArgs...)
	return tq.db.db.Exec(tq.string(), args...)
}

// Count .
func (tq *TableQuery) Count() (int64, error) {
	tq.action = sqlActionCount
	var count int64
	args := tq.setArgs
	args = append(args, tq.whereArgs...)
	row := tq.db.db.QueryRow(tq.string(), args...)
	err := row.Scan(&count)
	return count, err
}

func (tq *TableQuery) string() string {
	var q strings.Builder
	switch tq.action {
	case sqlActionDelete:
		q.WriteString("delete from ")
		q.WriteString(tq.table)
	case sqlActionCount:
		q.WriteString("select count(*) from ")
		q.WriteString(tq.table)
	case sqlActionUpdate:
		q.WriteString("update ")
		q.WriteString(tq.table)
	}
	for _, join := range tq.joins {
		q.WriteString(" join ")
		q.WriteString(join)
	}
	for i, set := range tq.sets {
		if i == 0 {
			q.WriteString(" set ")
		} else {
			q.WriteString(", ")
		}
		q.WriteString(set)
	}
	for i, where := range tq.wheres {
		if i == 0 {
			q.WriteString(" where ")
		} else {
			q.WriteString(where.conjunction)
		}
		q.WriteString(where.condition)
	}
	return q.String()
}
