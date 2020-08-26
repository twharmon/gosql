package gosql

import (
	"strings"
)

// CountQuery is a query for counting rows in a table.
type CountQuery struct {
	db         *DB
	count      string
	table      string
	joins      []string
	wheres     []*where
	whereArgs  []interface{}
	havings    []*having
	havingArgs []interface{}
	groupBy    string
}

// Where specifies which rows will be returned.
func (cq *CountQuery) Where(condition string, args ...interface{}) *CountQuery {
	w := &where{
		conjunction: " and ",
		condition:   condition,
	}
	cq.wheres = append(cq.wheres, w)
	cq.whereArgs = append(cq.whereArgs, args...)
	return cq
}

// OrWhere specifies which rows will be returned.
func (cq *CountQuery) OrWhere(condition string, args ...interface{}) *CountQuery {
	w := &where{
		conjunction: " or ",
		condition:   condition,
	}
	cq.wheres = append(cq.wheres, w)
	cq.whereArgs = append(cq.whereArgs, args...)
	return cq
}

// Join joins another table to this query.
func (cq *CountQuery) Join(join string) *CountQuery {
	cq.joins = append(cq.joins, join)
	return cq
}

// Exec executes the query.
func (cq *CountQuery) Exec() (int64, error) {
	var count int64
	row := cq.db.db.QueryRow(cq.String(), cq.whereArgs...)
	err := row.Scan(&count)
	return count, err
}

// String returns the string representation of CountQuery.
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
	if cq.groupBy != "" {
		q.WriteString(" group by ")
		q.WriteString(cq.groupBy)
	}
	for i, having := range cq.havings {
		if i == 0 {
			q.WriteString(" having ")
		} else {
			q.WriteString(having.conjunction)
		}
		q.WriteString(having.condition)
	}
	return q.String()
}

// Having specifies which rows will be returned.
func (cq *CountQuery) Having(condition string, args ...interface{}) *CountQuery {
	h := &having{
		conjunction: " and ",
		condition:   condition,
	}
	cq.havings = append(cq.havings, h)
	cq.havingArgs = append(cq.havingArgs, args...)
	return cq
}

// OrHaving specifies which rows will be returned.
func (cq *CountQuery) OrHaving(condition string, args ...interface{}) *CountQuery {
	h := &having{
		conjunction: " or ",
		condition:   condition,
	}
	cq.havings = append(cq.havings, h)
	cq.havingArgs = append(cq.havingArgs, args...)
	return cq
}

// GroupBy specifies how to group the results.
func (cq *CountQuery) GroupBy(bys ...string) *CountQuery {
	cq.groupBy = strings.Join(bys, ", ")
	return cq
}
