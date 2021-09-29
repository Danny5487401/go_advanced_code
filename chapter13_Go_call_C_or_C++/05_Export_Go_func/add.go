package main
//int GoAdd(int a, int b);
import "C"

//export GoAdd
func GoAdd(a,b C.int) C.int {
	return  a +b
}
