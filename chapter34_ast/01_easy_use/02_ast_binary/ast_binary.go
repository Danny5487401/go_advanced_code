package _2_ast_binary

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// Eval : 计算 expr 的值
func Eval(m map[string]int64, expr string) (bool, error) {
	// 解析表达式
	exprAst, err := parser.ParseExpr(expr)
	if err != nil {
		return false, err
	}

	// 打印 FileSet 形式
	fset := token.NewFileSet()
	ast.Print(fset, exprAst)

	return judge(exprAst, m), nil
}

// dfs 递归地计算表达式的值。
func judge(bop ast.Node, m map[string]int64) bool {
	// 叶子结点
	if isLeaf(bop) {
		// 断言成二元表达式
		expr := bop.(*ast.BinaryExpr)
		x := expr.X.(*ast.Ident)    // 左边
		y := expr.Y.(*ast.BasicLit) // 右边

		// 如果是 ">" 符号
		if expr.Op == token.GTR {
			left := m[x.Name]
			right, _ := strconv.ParseInt(y.Value, 10, 64)
			return left > right
		}
		return false
	}

	// 不是叶子节点那么一定是 binary expression（我们目前只处理二元表达式）
	expr, ok := bop.(*ast.BinaryExpr)
	if !ok {
		println("this cannot be true")
		return false
	}

	// 递归地计算左节点和右节点的值
	switch expr.Op {
	case token.LAND:
		return judge(expr.X, m) && judge(expr.Y, m)
	case token.LOR:
		return judge(expr.X, m) || judge(expr.Y, m)
	}

	println("unsupported operator")
	return false
}

// 判断是否是叶子节点
func isLeaf(bop ast.Node) bool {
	expr, ok := bop.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	// 二元表达式的最小单位，左节点是标识符，右节点是值
	_, okL := expr.X.(*ast.Ident)
	_, okR := expr.Y.(*ast.BasicLit)
	if okL && okR {
		return true
	}

	return false
}
