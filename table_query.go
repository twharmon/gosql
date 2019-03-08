package gosql

import (
	"strings"
)

const (
	sqlActionDelete = iota
	sqlActionCount
)

// TableQuery .
type TableQuery struct {
	db     *DB
	action int
	table  string
	joins  []string
	where  string
	args   []interface{}
}

// Table .
func (db *DB) Table(table string) *TableQuery {
	tq := new(TableQuery)
	tq.db = db
	tq.table = table
	return tq
}

// Where .
func (tq *TableQuery) Where(where string, args ...interface{}) *TableQuery {
	tq.where = where
	tq.args = args
	return tq
}

// Join .
func (tq *TableQuery) Join(join string) *TableQuery {
	tq.joins = append(tq.joins, join)
	return tq
}

// Delete .
func (tq *TableQuery) Delete() error {
	tq.action = sqlActionDelete
	_, err := tq.db.db.Exec(tq.string(), tq.args...)
	return err
}

// Count .
func (tq *TableQuery) Count() (int64, error) {
	tq.action = sqlActionCount
	var count int64
	row := tq.db.db.QueryRow(tq.string(), tq.args...)
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
	}
	for _, join := range tq.joins {
		q.WriteString(" join ")
		q.WriteString(join)
	}
	if tq.where != "" {
		q.WriteString(" where ")
		q.WriteString(tq.where)
	}
	return q.String()
}
