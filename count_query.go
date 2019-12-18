package gosql

import (
	"strings"
)

// CountQuery .
type CountQuery struct {
	db        *DB
	count     string
	table     string
	joins     []string
	wheres    []*where
	whereArgs []interface{}
}

// Where .
func (cq *CountQuery) Where(condition string, args ...interface{}) *CountQuery {
	w := &where{
		conjunction: " and ",
		condition:   condition,
	}
	cq.wheres = append(cq.wheres, w)
	cq.whereArgs = append(cq.whereArgs, args...)
	return cq
}

// OrWhere .
func (cq *CountQuery) OrWhere(condition string, args ...interface{}) *CountQuery {
	w := &where{
		conjunction: " or ",
		condition:   condition,
	}
	cq.wheres = append(cq.wheres, w)
	cq.whereArgs = append(cq.whereArgs, args...)
	return cq
}

// Join .
func (cq *CountQuery) Join(join string) *CountQuery {
	cq.joins = append(cq.joins, join)
	return cq
}

// Exec .
func (cq *CountQuery) Exec() (int64, error) {
	var count int64
	row := cq.db.db.QueryRow(cq.String(), cq.whereArgs...)
	err := row.Scan(&count)
	return count, err
}

// String .
func (cq *CountQuery) String() string {
	var q strings.Builder

	q.WriteString("select count(")
	q.WriteString(cq.count)
	q.WriteString(") from ")
	q.WriteString(cq.table)

	for _, join := range cq.joins {
		q.WriteString(" join ")
		q.WriteString(join)
	}
	for i, where := range cq.wheres {
		if i == 0 {
			q.WriteString(" where ")
		} else {
			q.WriteString(where.conjunction)
		}
		q.WriteString(where.condition)
	}
	return q.String()
}
