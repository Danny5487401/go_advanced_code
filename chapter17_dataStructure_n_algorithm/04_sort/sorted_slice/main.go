package main

import (
	"fmt"
	"go_advanced_code/chapter17_dataStructure_n_algorithm/04_sort/sortByReflect"
	"sort"
)

func main() {
	// 1。不同结构体比较使用
	type fruitList []interface{}
	list := fruitList{
		Fruit{"p苹果", 2, 3.3},
		Fruit{"x香蕉", 8, 4.55},
		Fruit{"j橘子", 5, 2.5},
		Fruit{"c橙子", 3, 6.05},
		User{1, "alice", true, 18},
	}

	sortByReflect.SortBodyByInt(list, "Name", "ASC")
	fmt.Println(list)

	//2.相同结构体使用
	users := []User{
		{1, "alice", true, 18},
		{2, "bob", false, 29},
		{3, "carl", true, 12},
		{4, "john", true, 16},
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Age < users[j].Age
	})
	fmt.Println(users)

	// 3.对字典的键名升序排序
	m := map[string]interface{}{"id": 1, "name": "admin", "pid": 0, "age": 18}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	//Go的sort包已经为基本数据类型都实现了sort功能，其函数名的最后一个字母是s，表示sort之意。比如：Ints, Float64s, Strings，等等。
	//对应的，int、float64、string等这些基本数据类型对应的集合类似，则被命名为IntSlice、Float64Slice、StringSlice等。
	//底层调用StringSlice结构体方法
	sort.Strings(keys)

	for i, v := range keys {
		fmt.Println(i, v)
	}

}

type Fruit struct {
	Name  string
	Count int
	Price float64
}

type User struct {
	Id     int
	Name   string
	Status bool
	Age    int
}
