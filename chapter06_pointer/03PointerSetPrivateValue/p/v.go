package p

import (
	"fmt"
)
/*
	利用unsafe包，可操作私有变量(在golang中称为“未导出变量”，变量名以小写字母开始)
 */
type Version struct {
	// 小写私有变量
	i int32
	j int64
}

func (v Version) PutI() {
	fmt.Printf("i=%d\n", v.i)
}

func (v Version) PutJ() {
	fmt.Printf("j=%d\n", v.j)
}
