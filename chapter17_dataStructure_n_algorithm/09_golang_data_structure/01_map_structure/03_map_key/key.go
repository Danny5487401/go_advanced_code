package main

import "fmt"

type _key struct {
}

type point struct {
	x int
	y int
}

type pair struct {
	x int
	y int
}

type Sumer interface {
	Sum() int
}

type Suber interface {
	Sub() int
}

func (p *pair) Sum() int {
	return p.x + p.y
}

func (p *point) Sum() int {
	return p.x + p.y
}

func (p pair) Sub() int {
	return p.x - p.y
}

func (p point) Sub() int {
	return p.x - p.y
}

func main() {
	fmt.Println("_key{} == _key{}: ", _key{} == _key{}) // output: true

	fmt.Println("point{} == point{}: ", point{x: 1, y: 2} == point{x: 1, y: 2})     // output: true
	fmt.Println("&point{} == &point{}: ", &point{x: 1, y: 2} == &point{x: 1, y: 2}) // output: false

	fmt.Println("[2]point{} == [2]point{}: ",
		[2]point{point{x: 1, y: 2}, point{x: 2, y: 3}} == [2]point{point{x: 1, y: 2}, point{x: 2, y: 3}}) //output: true

	var a Sumer = &pair{x: 1, y: 2}
	var a1 Sumer = &pair{x: 1, y: 2}
	var b Sumer = &point{x: 1, y: 2}
	fmt.Println("Sumer.byptr == Sumer.byptr: ", a == b)        // output: false
	fmt.Println("Sumer.sametype == Sumer.sametype: ", a == a1) // output: false

	var c Suber = pair{x: 1, y: 2}
	var d Suber = point{x: 1, y: 2}
	var d1 point = point{x: 1, y: 2}
	fmt.Println("Suber.byvalue == Suber.byvalue: ", c == d)  // output: false
	fmt.Println("Suber.byvalue == point.byvalue: ", d == d1) // output: true

	ci1 := make(chan int, 1)
	ci2 := ci1
	ci3 := make(chan int, 1)
	fmt.Println("chan int == chan int: ", ci1 == ci2) // output: true
	fmt.Println("chan int == chan int: ", ci1 == ci3) // output: false
}
