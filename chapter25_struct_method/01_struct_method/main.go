package main

type T struct {
	a int
}

func (t T) M1() {
	t.a = 10
}

func (t *T) M2() {
	t.a = 11
}

func main() {
	var t1 T
	t1.M1()
	t1.M2()

	var t2 = &T{}
	t2.M1()
	t2.M2()
}
