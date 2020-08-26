package gosql

import (
	"database/sql"
	"fmt"
	"strings"
)

// UpdateQuery holds information for an update query.
type UpdateQuery struct {
	db        *DB
	table     string
	joins     []string
	wheres    []*where
	sets      []string
	whereArgs []interface{}
	setArgs   []interface{}
}

// Where specifies which rows will be returned.
func (uq *UpdateQuery) Where(condition string, args ...interface{}) *UpdateQuery {
	w := &where{
		conjunction: " and ",
		condition:   condition,
	}
	uq.wheres = append(uq.wheres, w)
	uq.whereArgs = append(uq.whereArgs, args...)
	return uq
}

// OrWhere specifies which rows will be returned.
func (uq *UpdateQuery) OrWhere(condition string, args ...interface{}) *UpdateQuery {
	w := &where{
		conjunction: " or ",
		condition:   condition,
	}
	uq.wheres = append(uq.wheres, w)
	uq.whereArgs = append(uq.whereArgs, args...)
	return uq
}

// Set specifies how to update a row in a table.
func (uq *UpdateQuery) Set(set string, args ...interface{}) *UpdateQuery {
	uq.sets = append(uq.sets, set)
	uq.setArgs = append(uq.setArgs, args...)
	return uq
}

// Join joins another table to this query.
func (uq *UpdateQuery) Join(join string) *UpdateQuery {
	uq.joins = append(uq.joins, fmt.Sprintf(" join %s", join))
	return uq
}

// LeftJoin joins another table to this query.
func (uq *UpdateQuery) LeftJoin(join string) *UpdateQuery {
	uq.joins = append(uq.joins, fmt.Sprintf(" left join %s", join))
	return uq
}

// Exec executes the query.
func (uq *UpdateQuery) Exec() (sql.Result, error) {
	args := uq.setArgs
	args = append(args, uq.whereArgs...)
	return uq.db.db.Exec(uq.String(), args...)
}

// String returns the string representation of UpdateQuery.
func (uq *UpdateQuery) String() string {
	var q strings.Builder
	q.WriteString("update ")
	q.WriteString(uq.table)
	for _, join := range uq.joins {
		q.WriteString(join)
	}
	for i, set := range uq.sets {
		if i == 0 {
			q.WriteString(" set ")
		} else {
			q.WriteString(", ")
		}
		q.WriteString(set)
	}
	for i, where := range uq.wheres {
		if i == 0 {
			q.WriteString(" where ")
		} else {
			q.WriteString(where.conjunction)
		}
		q.WriteString(where.condition)
	}
	return q.String()
}
