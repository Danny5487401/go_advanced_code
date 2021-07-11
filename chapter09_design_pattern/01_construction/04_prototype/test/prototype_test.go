package test

import (
	"fmt"
	pt "go_advenced_code/chapter09_design_pattern/01_construction/04_prototype"
	"testing"
)

func Test_Prototype(t *testing.T) {
	u1 := pt.DefaultUserFactory.Create()
	fmt.Printf("u1 = %v\n", u1)

	u2 := pt.DefaultUserFactory.Create()
	fmt.Printf("u2 = %v\n", u2)
}
