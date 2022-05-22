package prototype

import (
	"fmt"
	"testing"
)

func Test_Prototype(t *testing.T) {
	u1 := DefaultUserFactory.Create()
	fmt.Printf("u1 = %v\n", u1)

	u2 := DefaultUserFactory.Create()
	fmt.Printf("u2 = %v\n", u2)
}
