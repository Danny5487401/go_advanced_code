package sp_fp_package

func Output(int) (int, int, int) // 汇编函数声明
/*
栈结构
	------
	ret2 (8 bytes)
	------
	ret1 (8 bytes)
	------
	ret0 (8 bytes)
	------
	arg0 (8 bytes)
	------ FP
	ret addr (8 bytes)
	------
	caller BP (8 bytes)
	------ pseudo SP
	frame content (8 bytes)
	------ hardware SP
*/
