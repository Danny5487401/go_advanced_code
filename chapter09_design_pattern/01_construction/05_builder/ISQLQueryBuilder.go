package builder

// 定义SQL查询表达式的建造者接口, 该接口定义了一系列步骤去创建复杂查询语句
type ISQLQueryBuilder interface {
	WithTable(table string) ISQLQueryBuilder //返回本身方便链式调用
	AddField(field string) ISQLQueryBuilder
	AddCondition(condition string) ISQLQueryBuilder
	WithOrderBy(orderBy string) ISQLQueryBuilder
	Build() ISQLQuery
}
