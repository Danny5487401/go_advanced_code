package main

type square struct {
	side int
}

func NewSquare(side int) shape {
	return &square{side: side}
}

func (s *square) accept(v visitor) {
	v.visitForSquare(s)
}

func (s *square) getType() string {
	return "Square"
}
