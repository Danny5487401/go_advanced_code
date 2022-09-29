package main

import (
	"fmt"
	oldFs "github.com/Danny5487401/go_advanced_code/chapter01_input_output/03_io_fs/01_go_1.15"
	newFs "github.com/Danny5487401/go_advanced_code/chapter01_input_output/03_io_fs/02_go_1.16"
)

func main() {
	// 打开文件
	fmt.Println(oldFs.FindExtFileGo115("chapter01_input_output/03_io_fs", ".txt"))
	fmt.Println(newFs.FindExtFileGo116("chapter01_input_output/03_io_fs", ".txt"))
}
