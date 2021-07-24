package builder

// 定义SQL查询表达式的接口
type ISQLQuery interface {
	ToSQL() string //生成sql语句
}
