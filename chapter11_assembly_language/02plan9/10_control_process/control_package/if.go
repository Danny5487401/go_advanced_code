package control_package

/*
if/goto跳转

三元表达式
	func If(ok bool, a, b int) int {
		if ok { return a } else { return b }
	}
汇编思维
	// 汇编语言中没有bool类型，我们改用int类型代替bool类型（真实的汇编是用byte表示bool类型，可以通过MOVBQZX指令加载byte类型的值，这里做了简化处理
	func If(ok int, a, b int) int {
		if ok == 0 { goto L }
		return a
	L:
		return b
	}
*/
