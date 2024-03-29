package builder

import (
	"fmt"
	"testing"
)

func Test_Builder(t *testing.T) {
	builder := NewSQLQueryBuilder()
	builder = builder.WithTable("product")
	builder = builder.AddField("id").AddField("name").AddField("price")
	builder = builder.AddCondition("enabled=1")
	builder = builder.WithOrderBy("price desc")
	query := builder.Build()
	fmt.Println(query.ToSQL())
}
