package gosql

import (
	"database/sql"
	"fmt"
	"strings"
)

// DeleteQuery is a query for deleting rows from a table.
type DeleteQuery struct {
	db        *DB
	execer    Execer
	table     string
	joins     []string
	wheres    []*where
	whereArgs []interface{}
}

// Where specifies which rows will be returned.
func (dq *DeleteQuery) Where(condition string, args ...interface{}) *DeleteQuery {
	w := &where{
		conjunction: " and ",
		condition:   condition,
	}
	dq.wheres = append(dq.wheres, w)
	dq.whereArgs = append(dq.whereArgs, args...)
	return dq
}

// OrWhere specifies which rows will be returned.
func (dq *DeleteQuery) OrWhere(condition string, args ...interface{}) *DeleteQuery {
	w := &where{
		conjunction: " or ",
		condition:   condition,
	}
	dq.wheres = append(dq.wheres, w)
	dq.whereArgs = append(dq.whereArgs, args...)
	return dq
}

// Join joins another table to this query.
func (dq *DeleteQuery) Join(join string) *DeleteQuery {
	dq.joins = append(dq.joins, fmt.Sprintf(" join %s", join))
	return dq
}

// LeftJoin joins another table to this query.
func (dq *DeleteQuery) LeftJoin(join string) *DeleteQuery {
	dq.joins = append(dq.joins, fmt.Sprintf(" left join %s", join))
	return dq
}

// Exec executes the query.
func (dq *DeleteQuery) Exec() (sql.Result, error) {
	return dq.execer.Exec(dq.String(), dq.whereArgs...)
}

// String returns the string representation of DeleteQuery.
func (dq *DeleteQuery) String() string {
	var q strings.Builder
	q.WriteString("delete from ")
	q.WriteString(dq.table)
	for _, join := range dq.joins {
		q.WriteString(join)
	}
	for i, where := range dq.wheres {
		if i == 0 {
			q.WriteString(" where ")
		} else {
			q.WriteString(where.conjunction)
		}
		q.WriteString(where.condition)
	}
	return q.String()
}
