package main

func main() {
	var m = make([]int, 10240)
	println(m[0])
}

/*
(⎈ |teleport.***)➜  03_escape_to_heap git:(feature/memory) ✗ go build -gcflags '-m'  04_big_slice.go
# command-line-arguments
./04_big_slice.go:3:6: can inline main
./04_big_slice.go:4:14: make([]int, 10240) escapes to heap

*/
