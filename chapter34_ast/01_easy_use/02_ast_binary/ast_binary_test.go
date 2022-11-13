package _2_ast_binary

import (
	"fmt"
	"testing"
)

func TestAst(t *testing.T) {

	m := map[string]int64{"orders": 100000, "driving_years": 18}

	// 单数超过 10000，且驾龄超过 5 年的老司机
	rule := `orders > 10000 && driving_years > 5`
	fmt.Println(Eval(m, rule))

}
