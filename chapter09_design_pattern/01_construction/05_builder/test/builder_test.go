package test

import (
	"fmt"
	bd "go_advenced_code/chapter09_design_pattern/01_construction/05_builder"
	"testing"
)

func Test_Builder(t *testing.T) {
	builder := bd.NewSQLQueryBuilder()
	builder = builder.WithTable("product")
	builder = builder.AddField("id").AddField("name").AddField("price")
	builder = builder.AddCondition("enabled=1")
	builder = builder.WithOrderBy("price desc")
	query := builder.Build()
	fmt.Println(query.ToSQL())
}
