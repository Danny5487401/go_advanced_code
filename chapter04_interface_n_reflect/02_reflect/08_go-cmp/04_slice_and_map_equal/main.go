package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func main() {
	s1 := []int{1, 2, 3, 4}
	s2 := []int{4, 3, 2, 1}
	fmt.Println("s1 equals s2?", cmp.Equal(s1, s2))                                                                       // 默认情况下，两个切片只有当长度相同，且对应位置上的元素都相等时，go-cmp才认为它们相等。
	fmt.Println("s1 equals s2 with option?", cmp.Equal(s1, s2, cmpopts.SortSlices(func(i, j int) bool { return i < j }))) // 排序好对比

	m1 := map[int]int{1: 10, 2: 20, 3: 30}
	m2 := map[int]int{1: 10, 2: 20, 3: 30}
	fmt.Println("m1 equals m2?", cmp.Equal(m1, m2))
	fmt.Println("m1 equals m2 with option?", cmp.Equal(m1, m2, cmpopts.SortMaps(func(i, j int) bool { return i < j }))) // cmpopts.SortMaps会将map[K]V类型按照键排序，生成一个[]struct{K, V}的切片，然后逐个比较

	u1 := []*User{
		{
			Name: "Joy",
			Age:  20,
		},
		{
			Name: "Danny",
			Age:  18,
		},
		{
			Name: "Jone",
			Age:  30,
		},
	}
	u2 := []*User{
		{
			Name: "Jone",
			Age:  30,
		},
		{
			Name: "Joy",
			Age:  20,
		},
		{
			Name: "Danny",
			Age:  18,
		},
	}

	fmt.Println("u1 equals u2?", cmp.Equal(u1, u2))                                                                                 // u1 equals u2? false
	fmt.Println("u1 equals u2 with option?", cmp.Equal(u1, u2, cmpopts.SortSlices(func(i, j *User) bool { return i.Age < j.Age }))) // u1 equals u2 with option? true
	fmt.Println("u1 equals u2 with option?", cmp.Equal(u1, u2, cmpopts.SortSlices(func(i, j User) bool { return i.Age < j.Age })))  // u1 equals u2 with option? false
}

type User struct {
	Name string
	Age  int
}
