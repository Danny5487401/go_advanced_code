package builder

import "strings"

//tSQLQuery实现了ISQLQuery接口, 根据各种参数生成复杂SQL语句
type tSQLQuery struct {
	table      string
	fields     []string
	conditions []string
	orderBy    string
}

func newSQLQuery() *tSQLQuery {
	return &tSQLQuery{
		table:      "",
		fields:     make([]string, 0),
		conditions: make([]string, 0),
		orderBy:    "",
	}
}

// Director的实现
// 生成sql语句
func (me *tSQLQuery) ToSQL() string {
	b := &strings.Builder{}
	b.WriteString("select ")

	for i, it := range me.fields {
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteString(it)
	}

	b.WriteString(" from ")
	b.WriteString(me.table)

	if len(me.conditions) > 0 {
		b.WriteString(" where ")
		for i, it := range me.conditions {
			if i > 0 {
				b.WriteString(" and ")
			}
			b.WriteString(it)
		}
	}

	if len(me.orderBy) > 0 {
		b.WriteString(" order by ")
		b.WriteString(me.orderBy)
	}

	return b.String()
}
