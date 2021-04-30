package main

//go:noinline
func stringParam(s string) {}

func main() {
	var x = "danny"
	stringParam(x)
}
