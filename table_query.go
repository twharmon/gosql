package gosql

import (
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
	wheres    []string
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
func (tq *TableQuery) Where(where string, arg interface{}) *TableQuery {
	tq.wheres = append(tq.wheres, where)
	tq.whereArgs = append(tq.whereArgs, arg)
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
func (tq *TableQuery) Delete() (int64, error) {
	tq.action = sqlActionDelete
	args := tq.setArgs
	args = append(args, tq.whereArgs...)
	res, err := tq.db.db.Exec(tq.string(), args...)
	if err != nil {
		return 0, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return ra, nil
}

// Update .
func (tq *TableQuery) Update() (int64, error) {
	tq.action = sqlActionUpdate
	args := tq.setArgs
	args = append(args, tq.whereArgs...)
	res, err := tq.db.db.Exec(tq.string(), args...)
	if err != nil {
		return 0, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return ra, nil
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
			q.WriteString(" and ")
		}
		q.WriteString(where)
	}
	return q.String()
}
