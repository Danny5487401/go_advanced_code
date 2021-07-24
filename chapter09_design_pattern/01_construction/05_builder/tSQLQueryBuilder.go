package builder

// tSQLQueryBuilder实现了ISQLQueryBuilder接口, 为各种参数设置提供了方法
type tSQLQueryBuilder struct {
	query *tSQLQuery //内部是复杂的信息
}

func NewSQLQueryBuilder() ISQLQueryBuilder {
	return &tSQLQueryBuilder{
		query: newSQLQuery(),
	}
}

func (me *tSQLQueryBuilder) WithTable(table string) ISQLQueryBuilder {
	me.query.table = table
	return me
}

func (me *tSQLQueryBuilder) AddField(field string) ISQLQueryBuilder {
	me.query.fields = append(me.query.fields, field)
	return me
}

func (me *tSQLQueryBuilder) AddCondition(condition string) ISQLQueryBuilder {
	me.query.conditions = append(me.query.conditions, condition)
	return me
}

func (me *tSQLQueryBuilder) WithOrderBy(orderBy string) ISQLQueryBuilder {
	me.query.orderBy = orderBy
	return me
}

// 最后一个Build方法用来返回我们创建的sql对象
func (me *tSQLQueryBuilder) Build() ISQLQuery {
	return me.query
}
