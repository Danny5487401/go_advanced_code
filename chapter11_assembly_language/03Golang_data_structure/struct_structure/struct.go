package main

type address struct {
	lng int
	lat int
}

type person struct {
	age    int
	height int
	addr   address
}

func readStruct(p person) (int, int, int, int) {
	return p.age,p.height,p.addr.lng,p.addr.lat
}

func main() {
	var p = person{
		age:    99,
		height: 88,
		addr: address{
			lng: 77,
			lat: 66,
		},
	}
	a, b, c, d := readStruct(p)
	println(a, b, c, d)
}
